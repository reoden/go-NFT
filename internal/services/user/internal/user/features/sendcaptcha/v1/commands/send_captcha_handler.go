package commands

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/random"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/sendcaptcha/v1/dtos"
)

type sendCaptchaHandler struct {
	fxparams.SendCaptchaHandlerParams
}

func NewSendCaptchaHandler(
	logger logger.Logger,
	userRepository contracts.UserRepository,
	cacheUserRepository contracts.UserCacheRepository,
	tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*SendCaptcha, *dtos.SendCaptchaResponseDto] {
	return &sendCaptchaHandler{
		SendCaptchaHandlerParams: fxparams.SendCaptchaHandlerParams{
			Log:             logger,
			UserRepository:  userRepository,
			RedisRepository: cacheUserRepository,
			Tracer:          tracer,
		},
	}
}

func (c *sendCaptchaHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*SendCaptcha, *dtos.SendCaptchaResponseDto](
		c,
	)
}

func (c *sendCaptchaHandler) Handle(
	ctx context.Context,
	command *SendCaptcha,
) (*dtos.SendCaptchaResponseDto, error) {
	captcha := random.String(6)
	err := c.RedisRepository.PutCaptcha(ctx, command.Phone, captcha)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("[Send_Captcha_Handler] get captcha telephone=%s from redis err=%+v", command.Phone, err),
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"[Send_Captcha_Handler] get captcha from redis = `%s`",
			captcha,
		),
		logger.Fields{"telephone": command.Phone},
	)

	phone := command.Phone

	var sendCaptchaResult *dtos.SendCaptchaResponseDto
	err = c.UserRepository.SendCaptcha(ctx, phone)
	if err != nil {
		return nil, err
	}

	sendCaptchaResult = &dtos.SendCaptchaResponseDto{}

	c.Log.Infow(
		fmt.Sprintf(
			"user with phone '%s' send Captcha = %s",
			command.Phone,
			captcha,
		),
		logger.Fields{
			"Telephone": command.Phone,
			"Captcha":   captcha,
		},
	)

	return sendCaptchaResult, err
}
