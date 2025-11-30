package user

import (
	"time"

	"github.com/reoden/go-NFT/pkg/testfixture"
	datamodel "github.com/reoden/go-NFT/user/internal/user/data/datamodels"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (ic *UserServiceConfigurator) seedUser(
	db *gorm.DB,
) error {
	err := seedDataManually(db)
	if err != nil {
		return err
	}

	return nil
}

func seedDataManually(gormDB *gorm.DB) error {
	var count int64

	// https://gorm.io/docs/advanced_query.html#Count
	gormDB.Model(&datamodel.UserDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	users := []*datamodel.UserDataModel{
		{
			UserId:    uuid.NewV4(),
			Nickname:  gofakeit.Name(),
			Phone:     gofakeit.Phone(),
			CreatedAt: time.Now(),
			//Description: gofakeit.AdjectiveDescriptive(),
			//Price:       gofakeit.Price(100, 1000),
		},
		{
			UserId:    uuid.NewV4(),
			Nickname:  gofakeit.Name(),
			Phone:     gofakeit.Phone(),
			CreatedAt: time.Now(),
			//Description: gofakeit.AdjectiveDescriptive(),
			//Price:       gofakeit.Price(100, 1000),
		},
	}

	err := gormDB.CreateInBatches(users, len(users)).Error
	if err != nil {
		return errors.Wrap(err, "error in seed database")
	}

	return nil
}

func seedDataWithFixture(gormDB *gorm.DB) error {
	var count int64

	// https://gorm.io/docs/advanced_query.html#Count
	gormDB.Model(&datamodel.UserDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	db, err := gormDB.DB()
	if err != nil {
		return errors.WrapIf(err, "error in seed database")
	}

	// https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	var data []struct {
		Nickname string
		UserId   uuid.UUID
		Phone    string
	}

	f := []struct {
		Nickname string
		UserId   uuid.UUID
		Phone    string
	}{
		{gofakeit.Name(), uuid.NewV4(), gofakeit.Phone()},
		{gofakeit.Name(), uuid.NewV4(), gofakeit.Phone()},
	}

	data = append(data, f...)

	err = testfixture.RunPostgresFixture(
		db,
		[]string{"db/fixtures/users"},
		map[string]interface{}{
			"Users": data,
		})
	if err != nil {
		return errors.WrapIf(err, "error in seed database")
	}

	return nil
}
