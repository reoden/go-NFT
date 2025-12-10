package commands

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/constants"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/logout/v1/dtos"
	"github.com/reoden/go-NFT/user/internal/user/models"
)

type logoutUserHandler struct {
	fxparams.LogoutHandlerParams
}

func NewLogoutUserHandler(
	logger logger.Logger,
	userRepository contracts.UserRepository,
	userOperateStreamRepository contracts.UserOperateStreamRepository,
	cacheUserRepository contracts.UserCacheRepository,
	tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*LogoutUser, *dtos.LogoutUserResponseDto] {
	return &logoutUserHandler{
		LogoutHandlerParams: fxparams.LogoutHandlerParams{
			Log:                         logger,
			UserRepository:              userRepository,
			UserOperateStreamRepository: userOperateStreamRepository,
			RedisRepository:             cacheUserRepository,
			Tracer:                      tracer,
		},
	}
}

func (c *logoutUserHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*LogoutUser, *dtos.LogoutUserResponseDto](
		c,
	)
}

func (c *logoutUserHandler) Handle(
	ctx context.Context,
	command *LogoutUser,
) (*dtos.LogoutUserResponseDto, error) {
	err := c.RedisRepository.AddTokenBlack(ctx, command.Token)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Logout_User_Handler] Add token=%s invalid to redis err=%+v", command.Token, err),
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"[Logout_User_Handler] make token = '%s' invalid successfully to redis",
			command.Token,
		),
		logger.Fields{"token": command.Token},
	)

	var logoutUserResult *dtos.LogoutUserResponseDto
	err = c.UserRepository.Logout(ctx, command.UserId)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UserId: command.UserId,
	}
	operateResult, err := c.UserOperateStreamRepository.InsertStream(ctx, user, constants.LOGOUT)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[Logout_User_Handler] insert stream err",
		)
	}
	c.Log.Infow(
		fmt.Sprintf(
			"user with user_id '%s' insert operate stream into database successfully",
			command.UserId.String(),
		),
		logger.Fields{
			"StreamId":    operateResult.Id,
			"UserId":      operateResult.UserId,
			"OperateType": operateResult.Type,
			"Param":       operateResult.Param,
		},
	)

	logoutUserResult = &dtos.LogoutUserResponseDto{}

	c.Log.Infow(
		fmt.Sprintf(
			"user '%s' with token = '%s' logout",
			command.UserId,
			command.Token,
		),
		logger.Fields{
			"Token":  command.Token,
			"UserId": command.UserId,
		},
	)

	return logoutUserResult, err
}
