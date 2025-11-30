package mediator

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	creatingUserCommondV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/commonds"
	createUserDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/dtos"
)

func ConfigUserMediator(
	logger logger.Logger,
	userDBContext *dbcontext.UserGormDBContext,
	userRepository contracts.UserRepository,
	cacheUserRepository contracts.UserCacheRepository,
	tracer tracing.AppTracer,
) error {
	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creatingUserCommondV1.CreateUser, *createUserDtosV1.CreateUserResponseDto](
		creatingUserCommondV1.NewCreateProductHandler(logger, userDBContext, userRepository, cacheUserRepository, tracer),
	)
	if err != nil {
		return err
	}

	//err = mediatr.RegisterRequestHandler[*getOrderByIdQueryV1.GetOrderById, *getOrderByIdDtosV1.GetOrderByIdResponseDto](
	//	getOrderByIdQueryV1.NewGetOrderByIdHandler(logger, mongoOrderReadRepository, tracer),
	//)
	//if err != nil {
	//	return err
	//}
	//
	//err = mediatr.RegisterRequestHandler[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](
	//	getOrdersQueryV1.NewGetOrdersHandler(logger, mongoOrderReadRepository, tracer),
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}
