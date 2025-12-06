package dtos

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// CreateUserRequestDto validation will handle in command level
type CreateUserRequestDto struct {
	Phone   string `json:"phone"`
	Captcha string `json:"captcha"`
}
