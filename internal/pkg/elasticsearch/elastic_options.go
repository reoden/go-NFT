package elasticsearch

import (
	"github.com/reoden/go-NFT/pkg/config"
	"github.com/reoden/go-NFT/pkg/config/environment"
	typeMapper "github.com/reoden/go-NFT/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[ElasticOptions]())

type ElasticOptions struct {
	URL string `mapstructure:"url"`
}

func provideConfig(environment environment.Environment) (*ElasticOptions, error) {
	return config.BindConfigKey[*ElasticOptions](optionName, environment)
}
