package dtos

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// AuthRequestDto validation will handle in command level
type AuthRequestDto struct {
	RealName string `json:"real_name"`
	IdCardNo string `json:"id_card_no"`
}
