package core

import (
	"github.com/reoden/go-echo-template/pkg/core/serializer/json"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"corefx",
	fx.Provide(
		json.NewDefaultJsonSerializer,
		json.NewDefaultEventJsonSerializer,
		json.NewDefaultMessageJsonSerializer,
		json.NewDefaultMetadataJsonSerializer,
	),
)
