package dbcontext

import (
	"github.com/reoden/go-NFT/pkg/postgresgorm/contracts"
	"github.com/reoden/go-NFT/pkg/postgresgorm/gormdbcontext"

	"gorm.io/gorm"
)

type UserGormDBContext struct {
	// our dbcontext base
	contracts.GormDBContext
}

func NewCatalogsDBContext(db *gorm.DB) *UserGormDBContext {
	// initialize base GormContext
	c := &UserGormDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
