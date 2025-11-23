package messagepersistence

import (
	"github.com/reoden/go-echo-template/pkg/postgresgorm/contracts"
	"github.com/reoden/go-echo-template/pkg/postgresgorm/gormdbcontext"

	"gorm.io/gorm"
)

type PostgresMessagePersistenceDBContext struct {
	// our dbcontext base
	contracts.GormDBContext
}

func NewPostgresMessagePersistenceDBContext(
	db *gorm.DB,
) *PostgresMessagePersistenceDBContext {
	// initialize base GormContext
	c := &PostgresMessagePersistenceDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
