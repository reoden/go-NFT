package projection

import (
	"context"

	"github.com/reoden/go-echo-template/pkg/es/models"
)

type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}
