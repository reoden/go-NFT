package contracts

import (
	"context"

	"github.com/reoden/go-NFT/user/internal/user/models"
)

type UserCacheRepository interface {
	PutUser(ctx context.Context, key string, product *models.User) error
	GetUserById(ctx context.Context, key string) (*models.User, error)
	DeleteUser(ctx context.Context, key string) error
	DeleteAllUsers(ctx context.Context) error
	PutCaptcha(ctx context.Context, key string, captcha string) error
	GetCaptcha(ctx context.Context, key string) (string, error)
}
