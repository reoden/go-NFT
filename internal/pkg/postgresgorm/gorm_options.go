package postgresgorm

import (
	"fmt"
	"path/filepath"

	"github.com/reoden/go-echo-template/pkg/config"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	typeMapper "github.com/reoden/go-echo-template/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GormOptions]())

type GormOptions struct {
	UseInMemory   bool   `mapstructure:"useInMemory"`
	UseSQLLite    bool   `mapstructure:"useSqlLite"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	DBName        string `mapstructure:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode"`
	Password      string `mapstructure:"password"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true"`
}

func (h *GormOptions) Dns() string {
	if h.UseInMemory {
		return ""
	}

	if h.UseSQLLite {
		projectRootDir := environment.GetProjectRootWorkingDirectory()
		dbFilePath := filepath.Join(projectRootDir, fmt.Sprintf("%s.db", h.DBName))

		return dbFilePath
	}

	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		h.User,
		h.Password,
		h.Host,
		h.Port,
		h.DBName,
	)

	return datasource
}

func provideConfig(environment environment.Environment) (*GormOptions, error) {
	cfg, err := config.BindConfigKey[*GormOptions](optionName, environment)
	fmt.Println("[DBEUG] cfg = ", cfg)
	return cfg, err
}
