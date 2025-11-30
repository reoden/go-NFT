package contracts

import (
	"context"

	"github.com/reoden/go-NFT/user/internal/user/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUserByTelephone(ctx context.Context, phone string) (*models.User, error)
}
