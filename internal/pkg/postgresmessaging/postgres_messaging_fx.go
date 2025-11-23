package postgresmessaging

import (
	"github.com/reoden/go-echo-template/pkg/core/messaging/persistmessage"
	"github.com/reoden/go-echo-template/pkg/postgresmessaging/messagepersistence"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"postgresmessagingfx",
	fx.Provide(
		messagepersistence.NewPostgresMessagePersistenceDBContext,
		messagepersistence.NewPostgresMessageService,
	),
	fx.Invoke(migrateMessaging),
)

func migrateMessaging(db *gorm.DB) error {
	err := db.Migrator().AutoMigrate(&persistmessage.StoreMessage{})

	return err
}
