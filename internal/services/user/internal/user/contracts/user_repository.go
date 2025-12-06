package contracts

import (
	"context"

	"github.com/reoden/go-NFT/user/internal/user/models"
	uuid "github.com/satori/go.uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
}
