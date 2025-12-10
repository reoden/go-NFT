package mediator

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/bloom"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
	creatingUserCommondV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/commonds"
	createUserDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/dtos"
	findUserByIdDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/finduserbyId/v1/dtos"
	findUserByIdQueryV1 "github.com/reoden/go-NFT/user/internal/user/features/finduserbyId/v1/queries"
	loginUserCommondV1 "github.com/reoden/go-NFT/user/internal/user/features/loginuser/v1/commands"
	loginUserDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/loginuser/v1/dtos"
	logoutCommondV1 "github.com/reoden/go-NFT/user/internal/user/features/logout/v1/commands"
	logoutDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/logout/v1/dtos"
	sendCaptchaCommondV1 "github.com/reoden/go-NFT/user/internal/user/features/sendcaptcha/v1/commands"
	sendCaptchaDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/sendcaptcha/v1/dtos"
)

func ConfigUserMediator(
	logger logger.Logger,
	userDBContext *dbcontext.UserGormDBContext,
	userRepository contracts.UserRepository,
	userOperateStreamRepository contracts.UserOperateStreamRepository,
	cacheUserRepository contracts.UserCacheRepository,
	bloomFilter *bloom.BloomFilterFactory,
	tracer tracing.AppTracer,
) error {
	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creatingUserCommondV1.CreateUser, *createUserDtosV1.CreateUserResponseDto](
		creatingUserCommondV1.NewCreateUserHandler(logger, userDBContext, userRepository, userOperateStreamRepository, cacheUserRepository, bloomFilter, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*findUserByIdQueryV1.FindUserById, *findUserByIdDtosV1.FindUserByIdResponseDto](
		findUserByIdQueryV1.NewFindUserByIdHandler(logger, userDBContext, userRepository, cacheUserRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*loginUserCommondV1.LoginUser, *loginUserDtosV1.LoginUserResponseDto](
		loginUserCommondV1.NewLoginUserHandler(logger, userDBContext, userRepository, userOperateStreamRepository, cacheUserRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*sendCaptchaCommondV1.SendCaptcha, *sendCaptchaDtosV1.SendCaptchaResponseDto](
		sendCaptchaCommondV1.NewSendCaptchaHandler(logger, userRepository, cacheUserRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*logoutCommondV1.LogoutUser, *logoutDtosV1.LogoutUserResponseDto](
		logoutCommondV1.NewLogoutUserHandler(logger, userRepository, userOperateStreamRepository, cacheUserRepository, tracer),
	)
	if err != nil {
		return err
	}
	//
	//err = mediatr.RegisterRequestHandler[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](
	//	getOrdersQueryV1.NewGetOrdersHandler(logger, mongoOrderReadRepository, tracer),
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}
