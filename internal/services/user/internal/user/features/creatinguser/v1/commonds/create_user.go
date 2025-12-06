package commonds

import (
	"regexp"
	"time"

	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type CreateUser struct {
	cqrs.Command
	UserId    uuid.UUID
	Captcha   string
	Phone     string
	CreatedAt time.Time
}

// NewCreateUser Create a new user
func NewCreateUser(
	phone string,
	captcha string,
) *CreateUser {
	command := &CreateUser{
		Command:   cqrs.NewCommandByT[CreateUser](),
		UserId:    uuid.NewV4(),
		Captcha:   captcha,
		Phone:     phone,
		CreatedAt: time.Now(),
	}

	return command
}

// NewCreateUserWithValidation Create a new user with inline validation - for defensive programming and ensuring validation even without using middleware
func NewCreateUserWithValidation(
	phone string,
	captcha string,
) (*CreateUser, error) {
	command := NewCreateUser(phone, captcha)
	err := command.Validate()

	return command, err
}

func (c *CreateUser) isTxRequest() {
}

func (c *CreateUser) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.UserId, validation.Required),
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
		validation.Field(&c.CreatedAt, validation.Required),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
