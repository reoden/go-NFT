package repositories

import (
	"context"
	"fmt"

	"github.com/reoden/go-NFT/pkg/core/data"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/otel/tracing/attribute"
	utils2 "github.com/reoden/go-NFT/pkg/otel/tracing/utils"
	"github.com/reoden/go-NFT/pkg/postgresgorm/repository"
	"github.com/reoden/go-NFT/user/internal/shared/constants"
	data2 "github.com/reoden/go-NFT/user/internal/user/contracts"
	"github.com/reoden/go-NFT/user/internal/user/models"
	uuid "github.com/satori/go.uuid"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	log                   logger.Logger
	gormGenericRepository data.GenericRepository[*models.User]
	tracer                tracing.AppTracer
}

func NewPostgresUserRepository(
	log logger.Logger,
	db *gorm.DB,
	tracer tracing.AppTracer,
) data2.UserRepository {
	gormRepository := repository.NewGenericGormRepository[*models.User](db)
	return &postgresUserRepository{
		log:                   log,
		gormGenericRepository: gormRepository,
		tracer:                tracer,
	}
}

func (p *postgresUserRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) (*models.User, error) {
	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.CreateUser")
	defer span.End()

	err := p.gormGenericRepository.Add(ctx, user)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the inserting user into the database.",
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("User", user))
	p.log.Infow(
		fmt.Sprintf(
			"user with user_id '%s' created",
			user.UserId,
		),
		logger.Fields{"User": user, "UserId": user.UserId, "Id": user.Id},
	)

	return user, nil
}

func (p *postgresUserRepository) FindUserById(
	ctx context.Context,
	userId uuid.UUID,
) (*models.User, error) {
	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.FindUserById")
	defer span.End()

	user, err := p.gormGenericRepository.FirstOrDefault(ctx, map[string]interface{}{
		"user_id": userId.String(),
	})
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the finding user into the database.",
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("User", user))
	p.log.Infow(
		fmt.Sprintf(
			"user with user_id '%s' found",
			userId.String(),
		),
		logger.Fields{
			"User":   user,
			"UserId": user.UserId,
			"Id":     user.Id,
			"Phone":  user.Phone,
		},
	)

	return user, nil
}

func (p *postgresUserRepository) UserLogin(ctx context.Context, telephone string) error {
	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.UserLogin")
	defer span.End()

	user, err := p.gormGenericRepository.FirstOrDefault(ctx, map[string]interface{}{
		"phone": telephone,
	})
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the finding user with telephone from the database.",
		),
	)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Object("User", user))
	p.log.Infow(
		fmt.Sprintf(
			"user with phone '%s' found",
			user.Phone,
		),
		logger.Fields{"User": user, "Phone": user.Phone},
	)

	return nil
}

func (p *postgresUserRepository) SendCaptcha(ctx context.Context, telephone string) error {
	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.SendCaptcha")
	defer span.End()

	p.log.Infow(
		fmt.Sprintf(
			"user with phone '%s' found",
			telephone,
		),
		logger.Fields{"Phone": telephone},
	)

	return nil
}

func (p *postgresUserRepository) Logout(ctx context.Context, userId uuid.UUID) error {
	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.Logout")
	defer span.End()

	p.log.Infow(
		fmt.Sprintf("user '%s' logout", userId.String()),
		logger.Fields{
			"user_id": userId.String(),
		},
	)
	return nil
}

func (p *postgresUserRepository) CheckAuth(ctx context.Context, userId uuid.UUID) (constants.UserStateEnum, error) {
	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.CheckAuth")
	defer span.End()

	user, err := p.gormGenericRepository.FirstOrDefault(ctx, map[string]interface{}{
		"user_id": userId,
	})

	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf("error in the finding user with user_id = '%s' from the database.", userId.String()),
		),
	)
	if err != nil {
		return "", err
	}

	span.SetAttributes(attribute.Object("User", user))
	p.log.Infow(
		fmt.Sprintf("user '%s' checked auth", userId.String()),
		logger.Fields{
			"user_id": userId.String(),
		},
	)

	return user.State, nil
}

//func (p *postgresUserRepository) GetAllUsers(
//	ctx context.Context,
//	listQuery *utils.ListQuery,
//) (*utils.ListResult[*models.User], error) {
//	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.GetAllUsers")
//	defer span.End()
//
//	result, err := p.gormGenericRepository.GetAll(ctx, listQuery)
//	err = utils2.TraceStatusFromContext(
//		ctx,
//		errors.WrapIf(
//			err,
//			"error in the paginate",
//		),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	p.log.Infow(
//		"users loaded",
//		logger.Fields{"UsersResult": result},
//	)
//
//	span.SetAttributes(attribute.Object("UsersResult", result))
//
//	return result, nil
//}
//
//func (p *postgresUserRepository) SearchUsers(
//	ctx context.Context,
//	searchText string,
//	listQuery *utils.ListQuery,
//) (*utils.ListResult[*models.User], error) {
//	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.SearchUsers")
//	span.SetAttributes(attribute2.String("SearchText", searchText))
//	defer span.End()
//
//	result, err := p.gormGenericRepository.Search(ctx, searchText, listQuery)
//	err = utils2.TraceStatusFromContext(
//		ctx,
//		errors.WrapIf(
//			err,
//			"error in the paginate",
//		),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	p.log.Infow(
//		fmt.Sprintf(
//			"users loaded for search term '%s'",
//			searchText,
//		),
//		logger.Fields{"UsersResult": result},
//	)
//	span.SetAttributes(attribute.Object("UsersResult", result))
//
//	return result, nil
//}
//
//func (p *postgresUserRepository) GetUserById(
//	ctx context.Context,
//	uuid uuid.UUID,
//) (*models.User, error) {
//	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.GetUserById")
//	span.SetAttributes(attribute2.String("Id", uuid.String()))
//	defer span.End()
//
//	user, err := p.gormGenericRepository.GetById(ctx, uuid)
//	err = utils2.TraceStatusFromSpan(
//		span,
//		errors.WrapIf(
//			err,
//			fmt.Sprintf(
//				"can't find the user with id %s into the database.",
//				uuid,
//			),
//		),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	span.SetAttributes(attribute.Object("User", user))
//	p.log.Infow(
//		fmt.Sprintf(
//			"user with id %s laoded",
//			uuid.String(),
//		),
//		logger.Fields{"User": user, "Id": uuid},
//	)
//
//	return user, nil
//}
//
//func (p *postgresUserRepository) UpdateUser(
//	ctx context.Context,
//	updateUser *models.User,
//) (*models.User, error) {
//	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.UpdateUser")
//	defer span.End()
//
//	err := p.gormGenericRepository.Update(ctx, updateUser)
//	err = utils2.TraceStatusFromSpan(
//		span,
//		errors.WrapIf(
//			err,
//			fmt.Sprintf(
//				"error in updating user with id %s into the database.",
//				updateUser.Id,
//			),
//		),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	span.SetAttributes(attribute.Object("User", updateUser))
//	p.log.Infow(
//		fmt.Sprintf(
//			"user with id '%s' updated",
//			updateUser.Id,
//		),
//		logger.Fields{
//			"User": updateUser,
//			"Id":   updateUser.Id,
//		},
//	)
//
//	return updateUser, nil
//}
//
//func (p *postgresUserRepository) DeleteUserByID(
//	ctx context.Context,
//	uuid uuid.UUID,
//) error {
//	ctx, span := p.tracer.Start(ctx, "postgresUserRepository.UpdateUser")
//	span.SetAttributes(attribute2.String("Id", uuid.String()))
//	defer span.End()
//
//	err := p.gormGenericRepository.Delete(ctx, uuid)
//	err = utils2.TraceStatusFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
//		"error in the deleting user with id %s into the database.",
//		uuid,
//	)))
//	if err != nil {
//		return err
//	}
//
//	p.log.Infow(
//		fmt.Sprintf(
//			"user with id %s deleted",
//			uuid,
//		),
//		logger.Fields{"User": uuid},
//	)
//
//	return nil
//}
