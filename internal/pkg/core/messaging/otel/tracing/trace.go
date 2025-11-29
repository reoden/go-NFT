package tracing

import (
	"github.com/reoden/go-NFT/pkg/otel/tracing"

	"go.opentelemetry.io/otel/trace"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/reoden/go-NFT/pkg/messaging",
	) // instrumentation name
}
