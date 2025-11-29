package domain

import (
	"github.com/reoden/go-NFT/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
