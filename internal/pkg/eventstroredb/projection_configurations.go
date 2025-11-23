package eventstroredb

import (
	"github.com/reoden/go-echo-template/pkg/es/contracts/projection"
)

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
