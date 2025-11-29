package dtos

import (
	"github.com/reoden/go-NFT/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `                      json:"listQuery"`
}
