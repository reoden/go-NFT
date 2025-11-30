package app

import (
    "github.com/reoden/go-NFT/pkg/fxapp"
    "github.com/reoden/go-NFT/pkg/fxapp/contracts"
)

type UserApplicationBuilder struct {
    contracts.ApplicationBuilder
}

func NewUserApplicationBuilder() *UserApplicationBuilder {
    builder := &UserApplicationBuilder{fxapp.NewApplicationBuilder()}

    return builder
}

func (a *UserApplicationBuilder) Build() *UserApplication {
    return NewUserApplication(
        a.GetProvides(),
        a.GetDecorates(),
        a.Options(),
        a.Logger(),
        a.Environment(),
    )
}
