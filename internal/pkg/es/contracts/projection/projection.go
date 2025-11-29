package projection

import (
	"context"

	"github.com/reoden/go-NFT/pkg/es/models"
)

type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
