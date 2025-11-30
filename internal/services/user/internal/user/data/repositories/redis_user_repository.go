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
	"github.com/reoden/go-NFT/user/internal/user/models"
	attribute2 "go.opentelemetry.io/otel/attribute"
)

const (
	redisUserMainPrefixKey       = "user_main_service"
	redisUserInviteCodePrefixKey = "user_invite_code_service_telephone"
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

func (r *redisUserRepository) GetInviteCode(ctx context.Context, key string) (string, error) {
	ctx, span := r.tracer.Start(ctx, "redisRepository.GetInviteCode")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserInviteCodePrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	redisKey := fmt.Sprintf("%s#%s", r.getRedisUserInviteCodePrefixKey(), key)
	inviteCodeBytes, err := r.redisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in getting inviteCode with Key %s from database",
					redisKey,
				),
			),
		)
	}

	inviteCode := string(inviteCodeBytes)

	span.SetAttributes(attribute.Object("invite_code", inviteCode))

	r.log.Infow(
		fmt.Sprintf(
			"invite_code with with key '%s', prefix '%s' laoded",
			redisKey,
			r.getRedisUserInviteCodePrefixKey(),
		),
		logger.Fields{
			"telephone":  key,
			"inviteCode": inviteCode,
			"Key":        redisKey,
			"PrefixKey":  r.getRedisUserInviteCodePrefixKey(),
		},
	)

	return inviteCode, nil
}

func (r *redisUserRepository) PutInviteCode(
	ctx context.Context,
	key string,
	inviteCode string,
) error {
	ctx, span := r.tracer.Start(ctx, "redisUserRepository.PutInviteCode")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisUserInviteCodePrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	redisKey := fmt.Sprintf("%s#%s", r.getRedisUserInviteCodePrefixKey(), key)
	if err := r.redisClient.SetNX(ctx, redisKey, inviteCode, time.Minute*time.Duration(5)).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating inviteCode with key %s",
					redisKey,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object(fmt.Sprintf("telephone#%s", key), inviteCode))

	r.log.Infow(
		fmt.Sprintf(
			"invite_code with key '%s', prefix '%s'  updated successfully",
			redisKey,
			r.getRedisUserMainPrefixKey(),
		),
		logger.Fields{
			"Phone":      key,
			"InviteCode": inviteCode,
			"Key":        redisKey,
			"PrefixKey":  r.getRedisUserMainPrefixKey(),
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

	if err := r.redisClient.HSetNX(ctx, r.getRedisUserMainPrefixKey(), key, userBytes).Err(); err != nil {
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
	//TODO implement me
	panic("implement me")
}

func (r *redisUserRepository) DeleteUser(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func (r *redisUserRepository) DeleteAllUsers(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *redisUserRepository) getRedisUserMainPrefixKey() string {
	return redisUserMainPrefixKey
}

func (r *redisUserRepository) getRedisUserInviteCodePrefixKey() string {
	return redisUserInviteCodePrefixKey
}
