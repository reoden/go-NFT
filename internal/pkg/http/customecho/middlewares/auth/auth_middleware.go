package auth

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// https://github.com/labstack/echo-jwt?tab=readme-ov-file#full-example

func EchoAuth(skipper func(c echo.Context) bool) echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		Skipper:    skipper,
	})
}
