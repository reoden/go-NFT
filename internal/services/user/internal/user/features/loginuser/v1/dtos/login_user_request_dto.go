package dtos

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// LoginUserRequestDto validation will handle in command level
type LoginUserRequestDto struct {
	Phone   string `json:"phone"`
	Captcha string `json:"captcha"`
}
