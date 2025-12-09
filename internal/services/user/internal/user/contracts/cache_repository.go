package contracts

import (
	"context"

	"github.com/reoden/go-NFT/user/internal/user/models"
)

type UserCacheRepository interface {
	PutUser(ctx context.Context, key string, product *models.User) error
	GetUserById(ctx context.Context, key string) (*models.User, error)
	PutCaptcha(ctx context.Context, key string, captcha string) error
	GetCaptcha(ctx context.Context, key string) (string, error)
	AddTokenBlack(ctx context.Context, token string) error
}
