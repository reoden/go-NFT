package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/redis/go-redis/v9"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/otel/tracing/attribute"
	"github.com/reoden/go-NFT/pkg/otel/tracing/utils"
	"github.com/reoden/go-NFT/user/internal/shared/constants"
	"github.com/reoden/go-NFT/user/internal/user/models"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

const (
	redisUserMainPrefixKey    = "user:cache:id:"
	redisUserCaptchaPrefixKey = "captcha:cache:"
)

type redisUserRepository struct {
	log         logger.Logger
	redisClient redis.UniversalClient
	tracer      tracing.AppTracer
}

func NewRedisUserRepository(
	log logger.Logger,
	redisClient redis.UniversalClient,
	tracer tracing.AppTracer,
) *redisUserRepository {
	return &redisUserRepository{
		log:         log,
		redisClient: redisClient,
		tracer:      tracer,
	}
}

func (r *redisUserRepository) GetCaptcha(ctx context.Context, key string) (string, error) {
	ctx, span := r.tracer.Start(ctx, "redisRepository.GetCaptcha")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserCaptchaPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	redisKey := fmt.Sprintf("%s%s", r.getRedisUserCaptchaPrefixKey(), key)
	captchaBytes, err := r.redisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in getting Captcha with Key %s from database",
					redisKey,
				),
			),
		)
	}

	captcha := string(captchaBytes)

	span.SetAttributes(attribute.Object("captcha", captcha))

	r.log.Infow(
		fmt.Sprintf(
			"captcha with with key '%s', prefix '%s' laoded",
			redisKey,
			r.getRedisUserCaptchaPrefixKey(),
		),
		logger.Fields{
			"telephone": key,
			"captcha":   captcha,
			"Key":       redisKey,
			"PrefixKey": r.getRedisUserCaptchaPrefixKey(),
		},
	)

	return captcha, nil
}

func (r *redisUserRepository) PutCaptcha(
	ctx context.Context,
	key string,
	captcha string,
) error {
	ctx, span := r.tracer.Start(ctx, "redisUserRepository.PutCaptcha")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserCaptchaPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	redisKey := fmt.Sprintf("%s%s", r.getRedisUserCaptchaPrefixKey(), key)
	if err := r.redisClient.SetNX(ctx, redisKey, captcha, constants.CaptchaExpireDuration).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating captcha with key %s",
					redisKey,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object(fmt.Sprintf("telephone#%s", key), captcha))

	r.log.Infow(
		fmt.Sprintf(
			"captcha with key '%s', prefix '%s'  updated successfully",
			redisKey,
			r.getRedisUserCaptchaPrefixKey(),
		),
		logger.Fields{
			"Phone":     key,
			"Captcha":   captcha,
			"Key":       redisKey,
			"PrefixKey": r.getRedisUserCaptchaPrefixKey(),
		},
	)

	return nil
}

func (r *redisUserRepository) PutUser(
	ctx context.Context,
	key string,
	user *models.User,
) error {
	ctx, span := r.tracer.Start(ctx, "redisUserRepository.PutUser")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserMainPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	cacheKey := fmt.Sprintf("%s%s", r.getRedisUserMainPrefixKey(), key)
	userBytes, err := json.Marshal(user)
	if err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error marshalling user",
			),
		)
	}

	if err := r.redisClient.SetNX(ctx, cacheKey, userBytes, constants.UserDataCacheExpireDuration).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating user with key %s",
					key,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("User", user))

	r.log.Infow(
		fmt.Sprintf(
			"user with key '%s', prefix '%s'  updated successfully",
			key,
			r.getRedisUserMainPrefixKey(),
		),
		logger.Fields{
			"User":      user,
			"UserId":    user.UserId,
			"Key":       key,
			"PrefixKey": r.getRedisUserMainPrefixKey(),
		},
	)

	return nil
}

func (r *redisUserRepository) GetUserById(ctx context.Context, key string) (*models.User, error) {
	ctx, span := r.tracer.Start(ctx, "redisRepository.GetUserById")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserMainPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	cacheKey := fmt.Sprintf("%s%s", r.getRedisUserMainPrefixKey(), key)

	userBytes, err := r.redisClient.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in getting User with Key %s from database",
					key,
				),
			),
		)
	}

	var user models.User
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return nil, utils.TraceErrStatusFromSpan(span, err)
	}

	span.SetAttributes(attribute.Object("User", user))

	r.log.Infow(
		fmt.Sprintf(
			"user with with key '%s', prefix '%s' laoded",
			key,
			r.getRedisUserMainPrefixKey(),
		),
		logger.Fields{
			"User":      user,
			"UserId":    user.UserId,
			"Key":       key,
			"PrefixKey": r.getRedisUserMainPrefixKey(),
		},
	)

	return &user, nil
}

func (r *redisUserRepository) AddTokenBlack(ctx context.Context, token string) error {
	ctx, span := r.tracer.Start(ctx, "redisUserRepository.AddTokenBlack")
	span.SetAttributes(
		attribute2.String("PrefixKey", constants.RedisTokenBlackPrefixKey),
	)

	key := fmt.Sprintf("%s%s", constants.RedisTokenBlackPrefixKey, token)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	if err := r.redisClient.SetNX(ctx, key, token, constants.UserTokenExpireDuration).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating token invalid with key %s",
					key,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Token", token))

	r.log.Infow(
		fmt.Sprintf(
			"token with key '%s', prefix '%s'  updated successfully",
			key,
			constants.RedisTokenBlackPrefixKey,
		),
		logger.Fields{
			"Token":     token,
			"Key":       key,
			"PrefixKey": constants.RedisTokenBlackPrefixKey,
		},
	)

	return nil
}

func (r *redisUserRepository) DelayedDelete(ctx context.Context, key string, delay time.Duration) error {
	ctx, span := r.tracer.Start(ctx, "redisRepository.DelayedDelete")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserMainPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	cacheKey := fmt.Sprintf("%s%s", r.getRedisUserMainPrefixKey(), key)

	go func() {
		select {
		case <-time.After(delay):
			r.redisClient.Del(ctx, cacheKey)
		case <-ctx.Done():
			return
		}
	}()
	return nil
}

func (r *redisUserRepository) DelUserById(ctx context.Context, key string) error {
	ctx, span := r.tracer.Start(ctx, "redisRepository.DelayedDelete")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserMainPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	cacheKey := fmt.Sprintf("%s%s", r.getRedisUserMainPrefixKey(), key)
	r.redisClient.Del(ctx, cacheKey)

	return nil
}

func (r *redisUserRepository) getRedisUserMainPrefixKey() string {
	return redisUserMainPrefixKey
}

func (r *redisUserRepository) getRedisUserCaptchaPrefixKey() string {
	return redisUserCaptchaPrefixKey
}
