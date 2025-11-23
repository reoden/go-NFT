package projection

import (
	"context"

	"github.com/reoden/go-echo-template/pkg/es/models"
)

type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
