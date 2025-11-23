package http

import (
	"github.com/reoden/go-echo-template/pkg/http/client"
	customEcho "github.com/reoden/go-echo-template/pkg/http/customecho"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("httpfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	client.Module,
	customEcho.Module,
)
