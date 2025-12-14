package queries

import (
	"context"
	"fmt"

	"github.com/reoden/go-NFT/pkg/core/cqrs"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/mapper"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/utils"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	dtosv1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/finduserbyId/v1/dtos"
	"github.com/reoden/go-NFT/user/internal/user/models"

	"github.com/mehdihadeli/go-mediatr"
)

type findUserByIdHandler struct {
	fxparams.FindUserByIdHandlerParams
}

func NewFindUserByIdHandler(
	logger logger.Logger,
	userDBContext *dbcontext.UserGormDBContext,
	userRepository contracts.UserRepository,
	cacheUserRepository contracts.UserCacheRepository,
	tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*FindUserById, *dtos.FindUserByIdResponseDto] {
	return &findUserByIdHandler{
		FindUserByIdHandlerParams: fxparams.FindUserByIdHandlerParams{
			Log:             logger,
			UserDBContext:   userDBContext,
			UserRepository:  userRepository,
			RedisRepository: cacheUserRepository,
			Tracer:          tracer,
		},
	}
}

func (c *findUserByIdHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*FindUserById, *dtos.FindUserByIdResponseDto](
		c,
	)
}

func (c *findUserByIdHandler) Handle(
	ctx context.Context,
	query *FindUserById,
) (*dtos.FindUserByIdResponseDto, error) {
	redisUser, err := c.RedisRepository.GetUserById(ctx, query.Id.String())

	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"error in getting user with id %d in the redis repository",
				query.Id,
			),
		)
	}

	var user *models.User

	if redisUser != nil {
		user = redisUser
	} else {
		var pgUser *models.User
		pgUser, err = c.UserRepository.FindUserById(ctx, query.Id)
		if err != nil {
			return nil, customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf("error in getting user with id %v in the postgres repository", query.Id),
			)
		}
		if pgUser == nil {
			pgUser, err = c.UserRepository.FindUserById(ctx, query.Id)
		}
		if err != nil {
			return nil, err
		}

		user = pgUser
		err = c.RedisRepository.PutUser(ctx, user.UserId.String(), user)
		if err != nil {
			return new(dtos.FindUserByIdResponseDto), err
		}
	}

	decodeRealName, err := utils.Decrypt(user.RealName)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in decrypting real name",
		)
	}
	user.RealName = decodeRealName
	decodeIdCardNo, err := utils.Decrypt(user.IdCardNo)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in decrypting id card no",
		)
	}
	user.IdCardNo = decodeIdCardNo
	userDto, err := mapper.Map[*dtosv1.UserDto](user)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping UserDto",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"user with userId '%v' found",
			query.Id,
		),
		logger.Fields{
			"UserId": query.Id,
		},
	)

	return &dtos.FindUserByIdResponseDto{User: userDto}, err
}
