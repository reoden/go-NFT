package commands

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/random"
	"github.com/reoden/go-NFT/pkg/bloom"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/mapper"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/postgresgorm/gormdbcontext"
	"github.com/reoden/go-NFT/user/internal/shared/constants"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	datamodel "github.com/reoden/go-NFT/user/internal/user/data/datamodels"
	dtosv1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/dtos"
	"github.com/reoden/go-NFT/user/internal/user/models"

	"github.com/mehdihadeli/go-mediatr"
)

type createUserHandler struct {
	fxparams.CreateUserHandlerParams

	nickNameBloomFilter   *bloom.BloomFilter
	inviteCodeBloomFilter *bloom.BloomFilter
}

func NewCreateUserHandler(
	logger logger.Logger,
	userDBContext *dbcontext.UserGormDBContext,
	userRepository contracts.UserRepository,
	userOperateStreamRepository contracts.UserOperateStreamRepository,
	cacheUserRepository contracts.UserCacheRepository,
	bloomFilter *bloom.BloomFilterFactory,
	tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*CreateUser, *dtos.CreateUserResponseDto] {
	return &createUserHandler{
		CreateUserHandlerParams: fxparams.CreateUserHandlerParams{
			Log:                         logger,
			UserDBContext:               userDBContext,
			UserRepository:              userRepository,
			UserOperateStreamRepository: userOperateStreamRepository,
			RedisRepository:             cacheUserRepository,
			BloomFilter:                 bloomFilter,
			Tracer:                      tracer,
		},
		nickNameBloomFilter:   bloomFilter.NewWithEstimates(1000000, 0.01, "nickname"),
		inviteCodeBloomFilter: bloomFilter.NewWithEstimates(1000000, 0.01, "inviteCode"),
	}
}

func (c *createUserHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*CreateUser, *dtos.CreateUserResponseDto](
		c,
	)
}

func (c *createUserHandler) Handle(
	ctx context.Context,
	command *CreateUser,
) (*dtos.CreateUserResponseDto, error) {
	userDataModelResult, err := gormdbcontext.FindDataModelByCond[*datamodel.UserDataModel](
		ctx,
		c.UserDBContext,
		map[string]any{
			"phone": command.Phone,
		},
	)

	if err != nil && !customErrors.IsNotFoundError(err) {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Create_User_Handler] find user with phone = %v err=%+v", command.Phone, err),
		)
	}

	if userDataModelResult != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			nil,
			fmt.Sprintf("[Create_User_Handler] user with phone = %v, this phone already exists", command.Phone),
		)
	}

	captcha, err := c.RedisRepository.GetCaptcha(ctx, command.Phone)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Create_User_Handler] get captcha telephone=%s from redis err=%+v", command.Phone, err),
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"[Create_User_Handler] get captcha from redis = `%s`",
			captcha,
		),
		logger.Fields{"telephone": command.Phone},
	)

	if captcha != command.Captcha {
		return nil, customErrors.NewApplicationErrorWrap(
			nil,
			fmt.Sprintf("[Create_User_Handler] get captcha=%s but need captcha = %s", command.Captcha, captcha),
		)
	}

	// generate nickname
	phone := command.Phone
	var (
		randomString    string
		defaultNickName string
	)
	for {
		randomString = random.String(6)
		defaultNickName = constants.DefaultNickNamePrefix + randomString + phone[7:11]
		ok, err := c.ExistsNickName(ctx, defaultNickName)
		if err != nil {
			return nil, customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf("[Create_User_Handler] check nickname=%s err=%+v", defaultNickName, err),
			)
		}

		if !ok {
			break
		}
	}

	user := &models.User{
		UserId:    command.UserId,
		Nickname:  defaultNickName,
		Phone:     command.Phone,
		CreatedAt: command.CreatedAt,
		State:     constants.User_INIT,
		UserRole:  constants.CUSTOMER,
	}

	var createUserResult *dtos.CreateUserResponseDto
	result, err := c.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	userDto, err := mapper.Map[*dtosv1.UserDto](result)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping UserDto",
		)
	}

	c.addNickname(ctx, userDto.Nickname)
	_ = c.RedisRepository.PutUser(ctx, userDto.UserId.String(), user)
	operateStreamResult, err := c.UserOperateStreamRepository.InsertStream(ctx, user, constants.REGISTER)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Create_User_Handler] insert user operate stream error"),
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"[Create_User_Handler] insert stream into user_operate_stream database = `%v`",
			operateStreamResult.Id,
		),
		logger.Fields{
			"StreamId":    operateStreamResult.Id,
			"UserId":      operateStreamResult.UserId,
			"OperateType": operateStreamResult.Type,
			"Param":       operateStreamResult.Param,
		},
	)

	createUserResult = &dtos.CreateUserResponseDto{
		UserID: user.UserId,
	}

	c.Log.Infow(
		fmt.Sprintf(
			"user with user_id '%s' created",
			command.UserId,
		),
		logger.Fields{
			"UserId": command.UserId,
		},
	)

	return createUserResult, err
}

func (c *createUserHandler) ExistsNickName(ctx context.Context, nickName string) (bool, error) {
	if c.nickNameBloomFilter != nil && c.nickNameBloomFilter.ExistsString(ctx, nickName) {
		_, err := gormdbcontext.FindDataModelByCond[*datamodel.UserDataModel](
			ctx,
			c.UserDBContext,
			map[string]any{
				"nickname": nickName,
			},
		)

		if err != nil && customErrors.IsNotFoundError(err) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (c *createUserHandler) addNickname(ctx context.Context, nickName string) {
	c.nickNameBloomFilter.AddString(ctx, nickName)
}
