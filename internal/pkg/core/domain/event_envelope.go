package domain

import (
	"github.com/reoden/go-echo-template/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
