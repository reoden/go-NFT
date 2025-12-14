package commands

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/mehdihadeli/go-mediatr"
    "github.com/reoden/go-NFT/pkg/authcertification"
    "github.com/reoden/go-NFT/pkg/core/cqrs"
    customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
    "github.com/reoden/go-NFT/pkg/logger"
    "github.com/reoden/go-NFT/pkg/otel/tracing"
    "github.com/reoden/go-NFT/pkg/postgresgorm/gormdbcontext"
    "github.com/reoden/go-NFT/pkg/utils"
    "github.com/reoden/go-NFT/user/internal/shared/constants"
    "github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
    "github.com/reoden/go-NFT/user/internal/user/contracts"
    datamodel "github.com/reoden/go-NFT/user/internal/user/data/datamodels"
    "github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
    "github.com/reoden/go-NFT/user/internal/user/features/checkauth/v1/dtos"
    "github.com/reoden/go-NFT/user/internal/user/models"
)

type authUserHandler struct {
    fxparams.CheckAuthHandlerParams
}

func NewAuthUserHandler(
    logger logger.Logger,
    dbContext *dbcontext.UserGormDBContext,
    userRepository contracts.UserRepository,
    userOperateStreamRepository contracts.UserOperateStreamRepository,
    cacheUserRepository contracts.UserCacheRepository,
    authRepo authcertification.AuthCertificationService,
    tracer tracing.AppTracer,
) cqrs.RequestHandlerWithRegisterer[*AuthUser, *dtos.AuthResponseDto] {
    return &authUserHandler{
        CheckAuthHandlerParams: fxparams.CheckAuthHandlerParams{
            Log:                         logger,
            UserDBContext:               dbContext,
            UserRepository:              userRepository,
            UserOperateStreamRepository: userOperateStreamRepository,
            RedisRepository:             cacheUserRepository,
            AuthRepo:                    authRepo,
            Tracer:                      tracer,
        },
    }
}

func (a *authUserHandler) RegisterHandler() error {
    return mediatr.RegisterRequestHandler[*AuthUser, *dtos.AuthResponseDto](
        a,
    )
}

func (a *authUserHandler) Handle(
    ctx context.Context,
    command *AuthUser,
) (*dtos.AuthResponseDto, error) {
    userModel, err := a.RedisRepository.GetUserById(ctx, command.UserId.String())
    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in GetUserById with user_id = '%v'", command.UserId),
        )
    }

    var userState constants.UserStateEnum

    if userModel != nil {
        userState = userModel.State
    } else {
        authState, err := a.UserRepository.CheckAuth(ctx, command.UserId)
        if err != nil {
            return nil, customErrors.NewApplicationErrorWrap(
                err,
                fmt.Sprintf("[authUserHandler.Handle] error in user_repository CheckAuth with user_id = '%v'", command.UserId),
            )
        }

        userState = authState
    }

    _ = a.RedisRepository.DelUserById(ctx, command.UserId.String())

    if userState == constants.User_AUTH ||
        userState == constants.User_ACTIVE {
        // 认证过
        return &dtos.AuthResponseDto{
            Msg:     "用户已认证",
            ErrCode: 0,
            Data:    userModel,
        }, nil
    }

    if userState != constants.User_INIT {
        return &dtos.AuthResponseDto{
            Msg:     "用户状态不能进行实名认证",
            ErrCode: http.StatusBadRequest,
            Data:    nil,
        }, nil
    }

    authRes, err := a.AuthRepo.Auth(ctx, command.RealName, command.IdCardNo)
    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in AuthRepo.Auth with real_name = '%v', idCardNo = '%v'", command.RealName, command.IdCardNo),
        )
    }

    if !authRes {
        return nil, customErrors.NewApplicationErrorWrap(
            nil,
            fmt.Sprintf("[authUserHandler.Handle] auth certification failed in AuthRepo.Auth with real_name = '%v', idCardNo = '%v'", command.RealName, command.IdCardNo),
        )
    }

    userDataModel, err := gormdbcontext.FindDataModelByCond[*datamodel.UserDataModel](
        ctx,
        a.UserDBContext,
        map[string]any{
            "user_id": command.UserId,
        },
    )

    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in FindDataModelByCond with user_id = '%v'", command.UserId),
        )
    }

    encodeRealName, err := utils.Encrypt(command.RealName)
    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in Encrypt with real_name = '%v'", command.RealName),
        )
    }

    encodeIdCardNo, err := utils.Encrypt(command.IdCardNo)
    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in Encrypt with id_card_no = '%v'", command.IdCardNo),
        )
    }
    userModel = &models.User{
        Id:            userDataModel.Id,
        UserId:        userDataModel.UserId,
        Nickname:      userDataModel.Nickname,
        Phone:         userDataModel.Phone,
        State:         constants.User_AUTH,
        Certification: true,
        RealName:      encodeRealName,
        IdCardNo:      encodeIdCardNo,
        UserRole:      userDataModel.UserRole,
        CreatedAt:     userDataModel.CreatedAt,
        UpdatedAt:     time.Now(),
    }

    updateUser, err := gormdbcontext.UpdateModel[*datamodel.UserDataModel, *models.User](
        ctx,
        a.UserDBContext,
        userModel,
    )

    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in UpdateModel with user_id = '%v'", command.UserId),
        )
    }

    operateStream, err := a.UserOperateStreamRepository.InsertStream(ctx, updateUser, constants.AUTH)
    if err != nil {
        return nil, customErrors.NewApplicationErrorWrap(
            err,
            fmt.Sprintf("[authUserHandler.Handle] error in InsertStream with user_id = '%v'", command.UserId),
        )
    }

    a.Log.Infow(
        fmt.Sprintf(
            "[Create_User_Handler] insert stream into user_operate_stream database = `%v`",
            operateStream.Id,
        ),
        logger.Fields{
            "StreamId":    operateStream.Id,
            "UserId":      operateStream.UserId,
            "OperateType": operateStream.Type,
            "Param":       operateStream.Param,
        },
    )

    authResponseResult := &dtos.AuthResponseDto{
        Msg:     "实名认证成功",
        ErrCode: 0,
        Data:    updateUser,
    }

    a.Log.Infow(
        fmt.Sprintf("user with id = '%v', auth certification successfully!", updateUser.UserId),
        logger.Fields{
            "UserId":   updateUser.UserId,
            "RealName": updateUser.RealName,
            "IdCardNo": updateUser.IdCardNo,
        },
    )

    return authResponseResult, nil
}
