package main

import (
	"context"
	"os"

	appconfig "github.com/reoden/go-echo-template/catalogs/config"
	"github.com/reoden/go-echo-template/pkg/config"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	"github.com/reoden/go-echo-template/pkg/logger"
	defaultLogger "github.com/reoden/go-echo-template/pkg/logger/defaultlogger"
	"github.com/reoden/go-echo-template/pkg/logger/external/fxlog"
	"github.com/reoden/go-echo-template/pkg/logger/zap"
	"github.com/reoden/go-echo-template/pkg/migration"
	"github.com/reoden/go-echo-template/pkg/migration/contracts"
	"github.com/reoden/go-echo-template/pkg/migration/goose"
	gormPostgres "github.com/reoden/go-echo-template/pkg/postgresgorm"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	// Add flags to specify the version
	cmdUp.Flags().Uint("version", 0, "Migration version")
	cmdDown.Flags().Uint("version", 0, "Migration version")

	// Add commands to the root command
	rootCmd.AddCommand(cmdUp)
	rootCmd.AddCommand(cmdDown)
}

var (
	rootCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "migration",
		Short: "A tool for running migrations",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute the "up" subcommand when no subcommand is specified
			if len(args) == 0 {
				cmd.SetArgs([]string{"up"})
				if err := cmd.Execute(); err != nil {
					defaultLogger.GetLogger().Error(err)
					os.Exit(1)
				}
			}
		},
	}

	cmdDown = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "down",
		Short: "Run a down migration",
		Run: func(cmd *cobra.Command, args []string) {
			executeMigration(cmd, migration.Down)
		},
	}

	cmdUp = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "up",
		Short: "Run an up migration",
		Run: func(cmd *cobra.Command, args []string) {
			executeMigration(cmd, migration.Up)
		},
	}
)

func executeMigration(cmd *cobra.Command, commandType migration.CommandType) {
	version, err := cmd.Flags().GetUint("version")
	if err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}

	app := fx.New(
		config.ModuleFunc(environment.Development),
		zap.Module,
		fxlog.FxLogger,
		gormPostgres.Module,
		appconfig.Module,
		//// use go-migrate library for migration
		//gomigrate.Module,
		// use go-migrate library for migration
		goose.Module,
		fx.Invoke(
			func(migrationRunner contracts.PostgresMigrationRunner, logger logger.Logger) {
				logger.Info("Migration process started...")
				switch commandType {
				case migration.Up:
					err = migrationRunner.Up(context.Background(), version)
				case migration.Down:
					err = migrationRunner.Down(context.Background(), version)
				}
				if err != nil {
					logger.Fatalf("migration failed, err: %s", err)
				}
				logger.Info("Migration completed...")
			},
		),
	)

	err = app.Start(context.Background())
	if err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}

	err = app.Stop(context.Background())
	if err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		defaultLogger.GetLogger().Error(err)
		os.Exit(1)
	}
}
