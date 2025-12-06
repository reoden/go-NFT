package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// FindUserByIdRequestDto validation will handle in command level
type FindUserByIdRequestDto struct {
	UserId uuid.UUID `param:"user_id" json:"user_id"`
}
