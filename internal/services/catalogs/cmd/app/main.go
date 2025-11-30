package main

import (
	"os"

	"github.com/reoden/go-NFT/catalogs/internal/shared/app"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "user-write-microservice",
	Short:            "user-write-microservice based on vertical slice architecture",
	Long:             `This is a command runner or cli for api architecture in golang.`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Reoden
// @contact.url https://github.com/Reoden
// @title Catalogs Service Api
// @version 1.0
// @description Catalogs Service Api.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Catalogs", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Write Service", pterm.FgLightMagenta.ToStyle())).
		Render()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
