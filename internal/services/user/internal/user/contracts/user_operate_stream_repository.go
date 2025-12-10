package contracts

import (
	"context"

	"github.com/reoden/go-NFT/user/internal/shared/constants"
	"github.com/reoden/go-NFT/user/internal/user/models"
)

type UserOperateStreamRepository interface {
	InsertStream(ctx context.Context, user *models.User, operateType constants.UserOperateTypeEnum) (*models.UserOperateStream, error)
}
