package commands

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-mediatr"
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
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/loginuser/v1/dtos"
	"github.com/reoden/go-NFT/user/internal/user/models"
)

type loginUserHandler struct {
	fxparams.UserLoginHandlerParams
}

func NewLoginUserHandler(
	logger logger.Logger,
	userDBContext *dbcontext.UserGormDBContext,
	userRepository contracts.UserRepository,
	userOperateStreamRepository contracts.UserOperateStreamRepository,
	cacheUserRepository contracts.UserCacheRepository,
	tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*LoginUser, *dtos.LoginUserResponseDto] {
	return &loginUserHandler{
		UserLoginHandlerParams: fxparams.UserLoginHandlerParams{
			Log:                         logger,
			UserDBContext:               userDBContext,
			UserRepository:              userRepository,
			UserOperateStreamRepository: userOperateStreamRepository,
			RedisRepository:             cacheUserRepository,
			Tracer:                      tracer,
		},
	}
}

func (c *loginUserHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*LoginUser, *dtos.LoginUserResponseDto](
		c,
	)
}

func (c *loginUserHandler) Handle(
	ctx context.Context,
	command *LoginUser,
) (*dtos.LoginUserResponseDto, error) {
	userDataModelResult, err := gormdbcontext.FindDataModelByCond[*datamodel.UserDataModel](
		ctx,
		c.UserDBContext,
		map[string]any{
			"phone": command.Phone,
		},
	)

	if err != nil && customErrors.IsNotFoundError(err) {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Login_User_Handler] find user with phone = %v err=%+v", command.Phone, err),
		)
	}

	if userDataModelResult == nil {
		return nil, customErrors.NewApplicationErrorWrap(
			nil,
			fmt.Sprintf("[Login_User_Handler] user with phone = %v, this phone does not exists", command.Phone),
		)
	}

	captcha, err := c.RedisRepository.GetCaptcha(ctx, command.Phone)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Login_User_Handler] get captcha telephone=%s from redis err=%+v", command.Phone, err),
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"[Login_User_Handler] get captcha from redis = `%s`",
			captcha,
		),
		logger.Fields{"telephone": command.Phone},
	)

	if captcha != command.Captcha {
		return nil, customErrors.NewApplicationErrorWrap(
			nil,
			fmt.Sprintf("[Login_User_Handler] get captcha=%s but need captcha = %s", command.Captcha, captcha),
		)
	}

	phone := command.Phone

	var loginUserResult *dtos.LoginUserResponseDto
	err = c.UserRepository.UserLogin(ctx, phone)
	if err != nil {
		return nil, err
	}

	user, err := mapper.Map[*models.User](userDataModelResult)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping model user",
		)
	}

	operateResult, err := c.UserOperateStreamRepository.InsertStream(ctx, user, constants.LOGIN)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[Login_User_Handler] insert stream err",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"user with phone '%s' insert operate stream into database successfully",
			command.Phone,
		),
		logger.Fields{
			"StreamId":    operateResult.Id,
			"UserId":      operateResult.UserId,
			"OperateType": operateResult.Type,
			"Param":       operateResult.Param,
		},
	)

	loginUserResult = &dtos.LoginUserResponseDto{
		UserId: userDataModelResult.UserId,
	}

	c.Log.Infow(
		fmt.Sprintf(
			"user with phone '%s' login",
			command.Phone,
		),
		logger.Fields{
			"Telephone": command.Phone,
		},
	)

	return loginUserResult, err
}
