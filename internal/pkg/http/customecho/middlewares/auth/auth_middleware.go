package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// https://github.com/labstack/echo-jwt?tab=readme-ov-file#full-example

type TokenBlacklistChecker interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

func EchoAuth(skipper func(c echo.Context) bool) echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		Skipper:    skipper,
	})
}

func JWTWithBlacklist(
	jwtMiddleware echo.MiddlewareFunc,
	checker TokenBlacklistChecker,
	skipper func(c echo.Context) bool,
) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper != nil && skipper(c) {
				return next(c)
			}
			if err := jwtMiddleware(func(c echo.Context) error { return nil })(c); err != nil {
				return err
			}

			token := extractRawToken(c)
			if token == "" {
				return echo.ErrUnauthorized
			}

			blacklisted, err := checker.IsBlacklisted(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "blacklist check error")
			}
			if blacklisted {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

func extractRawToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
