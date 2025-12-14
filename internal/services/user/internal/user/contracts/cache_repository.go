package contracts

import (
	"context"
	"time"

	"github.com/reoden/go-NFT/user/internal/user/models"
)

type UserCacheRepository interface {
	PutUser(ctx context.Context, key string, user *models.User) error
	GetUserById(ctx context.Context, key string) (*models.User, error)
	PutCaptcha(ctx context.Context, key string, captcha string) error
	GetCaptcha(ctx context.Context, key string) (string, error)
	AddTokenBlack(ctx context.Context, token string) error
	DelayedDelete(ctx context.Context, key string, delay time.Duration) error
	DelUserById(ctx context.Context, key string) error
}
