package commands

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator
type SendCaptcha struct {
	cqrs.Command
	Phone string `json:"phone"`
}

// NewSendCaptcha user send captcha
func NewSendCaptcha(
	phone string,
) *SendCaptcha {
	command := &SendCaptcha{
		Command: cqrs.NewCommandByT[SendCaptcha](),
		Phone:   phone,
	}

	return command
}

// NewSendCaptchaWithValidation user send captcha with inline validation - for defensive programming and ensuring validation even without using middleware
func NewSendCaptchaWithValidation(
	phone string,
) (*SendCaptcha, error) {
	command := NewSendCaptcha(phone)
	err := command.Validate()

	return command, err
}

func (c *SendCaptcha) isTxRequest() {
}

func (c *SendCaptcha) Validate() error {
	err := validation.ValidateStruct(
		c,
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
