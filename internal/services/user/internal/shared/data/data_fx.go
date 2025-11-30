package data

import (
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"

	"go.uber.org/fx"
)

// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"datafx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		dbcontext.NewUserDBContext,
	),
)
