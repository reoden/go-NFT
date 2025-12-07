package queries

import (
	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	validation "github.com/go-ozzo/ozzo-validation"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type FindUserById struct {
	cqrs.Query
	Id uuid.UUID
}

// NewFindUserById find a user
func NewFindUserById(
	id uuid.UUID,
) *FindUserById {
	command := &FindUserById{
		Query: cqrs.NewQueryByT[FindUserById](),
		Id:    id,
	}

	return command
}

// NewFindUserByIdWithValidation find a user with inline validation - for defensive programming and ensuring validation even without using middleware
func NewFindUserByIdWithValidation(
	id uuid.UUID,
) (*FindUserById, error) {
	command := NewFindUserById(id)
	err := command.Validate()

	return command, err
}

func (c *FindUserById) isTxRequest() {
}

func (c *FindUserById) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Id, validation.Required),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
