package authcertification

import (
    "github.com/iancoleman/strcase"
    "github.com/reoden/go-NFT/pkg/config"
    "github.com/reoden/go-NFT/pkg/config/environment"
    typeMapper "github.com/reoden/go-NFT/pkg/reflection/typemapper"
)

type AuthCertificationOptions struct {
    Host    string `mapstructure:"host"`
    Path    string `mapstructure:"path"`
    AppCode string `mapstructure:"appcode"`
}

func provideConfig(
    environment environment.Environment,
) (*AuthCertificationOptions, error) {
    optionName := strcase.ToLowerCamel(
        typeMapper.GetGenericTypeNameByT[AuthCertificationOptions](),
    )
    return config.BindConfigKey[*AuthCertificationOptions](optionName, environment)
}