package dtos

import "github.com/reoden/go-NFT/pkg/core/serializer/json"

type AuthResponseDto struct {
	ErrCode int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}

func (a *AuthResponseDto) String() string {
	return json.PrettyPrint(a)
}
