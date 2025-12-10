package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/reoden/go-NFT/pkg/core/data"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/otel/tracing/attribute"
	utils2 "github.com/reoden/go-NFT/pkg/otel/tracing/utils"
	"github.com/reoden/go-NFT/pkg/postgresgorm/repository"
	"github.com/reoden/go-NFT/user/internal/shared/constants"
	data2 "github.com/reoden/go-NFT/user/internal/user/contracts"
	"github.com/reoden/go-NFT/user/internal/user/models"
	"gorm.io/gorm"
)

type postgresUserOperateStreamRepository struct {
	log                   logger.Logger
	gormGenericRepository data.GenericRepository[*models.UserOperateStream]
	tracer                tracing.AppTracer
}

func NewPostgresUserOperateStreamRepository(
	log logger.Logger,
	db *gorm.DB,
	tracer tracing.AppTracer,
) data2.UserOperateStreamRepository {
	gormRepository := repository.NewGenericGormRepository[*models.UserOperateStream](db)
	return &postgresUserOperateStreamRepository{
		log:                   log,
		gormGenericRepository: gormRepository,
		tracer:                tracer,
	}
}

func (p *postgresUserOperateStreamRepository) InsertStream(
	ctx context.Context,
	user *models.User,
	operateType constants.UserOperateTypeEnum,
) (*models.UserOperateStream, error) {
	ctx, span := p.tracer.Start(ctx, "postgresUserOperateStreamRepository.InsertStream")
	defer span.End()

	userOperateStream := &models.UserOperateStream{
		UserId:      user.UserId,
		Type:        string(operateType),
		OperateTime: time.Now(),
		GMTCreate:   time.Now(),
		GMTModified: time.Now(),
	}

	userBytes, err := json.Marshal(user)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the marshaling user into json.",
		),
	)
	if err != nil {
		return nil, err
	}

	userOperateStream.Param = string(userBytes)

	err = p.gormGenericRepository.Add(ctx, userOperateStream)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the inserting user operate stream into the database.",
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("UserOperateStream", userOperateStream))
	p.log.Infow(
		fmt.Sprintf(
			"user operate stream with user_id '%s' created",
			user.UserId.String(),
		),
		logger.Fields{"UserOperateStream": userOperateStream, "UserId": user.UserId.String(), "Id": userOperateStream.Id},
	)

	return userOperateStream, nil
}
