package auth

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/reoden/go-NFT/pkg/constants"
)

type RedisTokenBlacklistChecker struct {
	client redis.UniversalClient
}

func NewRedisTokenBlacklistChecker(client *redis.Client) TokenBlacklistChecker {
	return &RedisTokenBlacklistChecker{client: client}
}

func (r *RedisTokenBlacklistChecker) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("%s#%s", constants.RedisTokenBlackPrefixKey, token)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
