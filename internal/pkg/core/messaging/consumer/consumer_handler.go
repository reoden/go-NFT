package consumer

import (
	"context"

	"github.com/reoden/go-NFT/pkg/core/messaging/types"
)

type ConsumerHandler interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
