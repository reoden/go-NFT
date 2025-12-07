package commands

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type LoginUser struct {
	cqrs.Command
	Captcha string
	Phone   string
}

// NewLoginUser user loginuser
func NewLoginUser(
	phone string,
	captcha string,
) *LoginUser {
	command := &LoginUser{
		Command: cqrs.NewCommandByT[LoginUser](),
		Captcha: captcha,
		Phone:   phone,
	}

	return command
}

// NewLoginUserWithValidation user loginuser with inline validation - for defensive programming and ensuring validation even without using middleware
func NewLoginUserWithValidation(
	phone string,
	captcha string,
) (*LoginUser, error) {
	command := NewLoginUser(phone, captcha)
	err := command.Validate()

	return command, err
}

func (c *LoginUser) isTxRequest() {
}

func (c *LoginUser) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(
			&c.Captcha,
			validation.Required,
			validation.Length(0, 6),
		),
		validation.Field(
			&c.Phone,
			validation.Required,
			validation.Match(regexp.MustCompile(`^1[3-9]\d{9}$`)),
			validation.Length(0, 11),
		),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
