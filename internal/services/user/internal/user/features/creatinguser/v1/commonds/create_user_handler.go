package commonds

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/random"
	"github.com/reoden/go-NFT/pkg/constants"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/mapper"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/postgresgorm/gormdbcontext"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	datamodel "github.com/reoden/go-NFT/user/internal/user/data/datamodels"
	dtosv1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/dtos"
	"github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/events/integrationevents"
	"github.com/reoden/go-NFT/user/internal/user/models"

	"github.com/mehdihadeli/go-mediatr"
)

type createUserHandler struct {
	fxparams.CreateUserHandlerParams
}

func NewCreateProductHandler(
	logger logger.Logger,
	userDBContext *dbcontext.UserGormDBContext,
	userRepository contracts.UserRepository,
	cacheUserRepository contracts.UserCacheRepository,
	tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*CreateUser, *dtos.CreateUserResponseDto] {
	return &createUserHandler{
		CreateUserHandlerParams: fxparams.CreateUserHandlerParams{
			Log:             logger,
			UserDBContext:   userDBContext,
			UserRepository:  userRepository,
			RedisRepository: cacheUserRepository,
			Tracer:          tracer,
		},
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

	inviteCode, err := c.RedisRepository.GetInviteCode(ctx, command.Phone)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Create_User_Handler] get invite_code telephone=%s from redis err=%+v", command.Phone, err),
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"[Create_User_Handler] get invite_code from redis = `%s`",
			inviteCode,
		),
		logger.Fields{"telephone": command.Phone},
	)

	if inviteCode != command.InviteCode {
		return nil, customErrors.NewApplicationErrorWrap(
			nil,
			fmt.Sprintf("[Create_User_Handler] get invite_code=%s but need invite_code = %s", command.InviteCode, inviteCode),
		)
	}

	phone := command.Phone
	randomString := random.String(6)
	defaultNickName := constants.DEFAULT_NICK_NAME_PREFIX + randomString + phone[7:11]
	user := &models.User{
		UserId:    command.UserId,
		Nickname:  defaultNickName,
		Phone:     command.Phone,
		CreatedAt: command.CreatedAt,
	}

	var createProductResult *dtos.CreateUserResponseDto

	//result, err := gormdbcontext.AddModel[*datamodel.UserDataModel, *models.User](
	//	ctx,
	//	c.UserDBContext,
	//	user,
	//)
	result, err := c.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	userDto, err := mapper.Map[*dtosv1.UserDto](result)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ProductDto",
		)
	}

	userCreated := integrationevents.NewUserCreatedV1(
		userDto,
	)

	c.Log.Infow(
		fmt.Sprintf(
			"UserCreated message with messageId `%s` published to the rabbitmq broker",
			userCreated.MessageId,
		),
		logger.Fields{"MessageId": userCreated.MessageId},
	)

	createProductResult = &dtos.CreateUserResponseDto{
		UserID: user.UserId,
	}

	c.Log.Infow(
		fmt.Sprintf(
			"user with user_id '%s' created",
			command.UserId,
		),
		logger.Fields{
			"UserId":    command.UserId,
			"MessageId": userCreated.MessageId,
		},
	)

	return createProductResult, err
}
