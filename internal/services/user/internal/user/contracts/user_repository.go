package contracts

import (
	"context"

	"github.com/reoden/go-NFT/user/internal/user/models"
	uuid "github.com/satori/go.uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	UserLogin(ctx context.Context, telephone string) error
	SendCaptcha(ctx context.Context, telephone string) error
	Logout(ctx context.Context, userId uuid.UUID) error
}
