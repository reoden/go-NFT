package commands

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type LogoutUser struct {
	cqrs.Command
	Token  string
	UserId uuid.UUID
}

func NewLogoutUser(
	token string,
	userId uuid.UUID,
) *LogoutUser {
	command := &LogoutUser{
		Command: cqrs.NewCommandByT[LogoutUser](),
		Token:   token,
		UserId:  userId,
	}

	return command
}

func NewLogoutUserWithValidation(
	token string,
	userId uuid.UUID,
) (*LogoutUser, error) {
	command := NewLogoutUser(token, userId)
	err := command.Validate()

	return command, err
}

func (l *LogoutUser) isTxRequest() {
}
func (l *LogoutUser) Validate() error {
	err := validation.ValidateStruct(
		l,
		validation.Field(
			&l.Token,
			validation.Required,
		),
		validation.Field(
			&l.UserId,
			validation.Required,
		),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
