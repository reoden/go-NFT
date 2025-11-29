package app

import (
	"github.com/reoden/go-NFT/pkg/fxapp"
	"github.com/reoden/go-NFT/pkg/fxapp/contracts"
)

type CatalogsApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewCatalogsWriteApplicationBuilder() *CatalogsApplicationBuilder {
	builder := &CatalogsApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (a *CatalogsApplicationBuilder) Build() *CatalogsApplication {
	return NewCatalogsApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
