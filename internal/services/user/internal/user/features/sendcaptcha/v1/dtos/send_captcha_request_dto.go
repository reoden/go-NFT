package dtos

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// SendCaptchaRequestDto validation will handle in command level
type SendCaptchaRequestDto struct {
	Phone string `json:"phone"`
}
