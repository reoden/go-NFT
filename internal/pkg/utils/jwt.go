package utils

import (
	"os"
	"time"

	"emperror.dev/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/reoden/go-NFT/pkg/constants"
	uuid "github.com/satori/go.uuid"
)

func GenJWTToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId.String(),
		"exp":    time.Now().Add(constants.TokenExpireDuration).Unix(),
		"iat":    time.Now().Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := rawToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseJWTToken(c echo.Context) (string, uuid.UUID, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return "", uuid.Nil, errors.New(constants.ErrJWTTokenInvalid)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", uuid.Nil, errors.New(constants.ErrJWTTokenFailedCastClaim)
	}
	uuidString := claims["userId"].(string)
	userId, err := uuid.FromString(uuidString)
	return token.Raw, userId, err
}
