package commands

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator
type AuthUser struct {
	cqrs.Command
	RealName string
	IdCardNo string
	UserId   uuid.UUID
}

func NewAuthUser(
	realName string,
	idCardNo string,
	userId uuid.UUID,
) *AuthUser {
	command := &AuthUser{
		RealName: realName,
		IdCardNo: idCardNo,
		UserId:   userId,
	}

	return command
}

func NewAuthUserWithValidation(
	realName string,
	idCardNo string,
	userId uuid.UUID,
) (*AuthUser, error) {
	command := NewAuthUser(realName, idCardNo, userId)
	if err := command.Validate(); err != nil {
		return nil, err
	}

	return command, nil
}

func (c *AuthUser) isTxRequest() {
}

func (c *AuthUser) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.RealName, validation.Required),
		validation.Field(
			&c.IdCardNo,
			validation.Required,
			validation.Length(0, 18),
		),
		validation.Field(&c.UserId, validation.Required),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
