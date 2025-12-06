package grpc

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing/attribute"
	"github.com/reoden/go-NFT/user/internal/shared/contracts"
	userService "github.com/reoden/go-NFT/user/internal/shared/grpc/genproto"
	createUserCommandV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/commonds"
	createUserDtosV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/dtos"
	attribute2 "go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var grpcMetricsAttr = api.WithAttributes(
	attribute2.Key("MetricsType").String("Http"),
)

type UserGrpcServiceServer struct {
	userMetrics *contracts.UserMetrics
	logger      logger.Logger
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	// product_service_client.UnimplementedProductsServiceServer
}

func NewUserGrpcService(
	userMetrics *contracts.UserMetrics,
	logger logger.Logger,
) *UserGrpcServiceServer {
	return &UserGrpcServiceServer{
		userMetrics: userMetrics,
		logger:      logger,
	}
}

func (s *UserGrpcServiceServer) CreateUser(
	ctx context.Context,
	req *userService.CreateUserReq,
) (*userService.CreateUserRes, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))
	s.userMetrics.CreateUserGrpcRequests.Add(ctx, 1, grpcMetricsAttr)

	command, err := createUserCommandV1.NewCreateUserWithValidation(
		req.GetPhone(),
		req.GetCaptcha(),
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[UserGrpcServiceServer_CreateUser.StructCtx] command validation failed",
		)
		s.logger.Errorf(
			fmt.Sprintf(
				"[UserGrpcServiceServer_CreateUser.StructCtx] err: %v",
				validationErr,
			),
		)
		return nil, validationErr
	}

	result, err := mediatr.Send[*createUserCommandV1.CreateUser, *createUserDtosV1.CreateUserResponseDto](
		ctx,
		command,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[UserGrpcServiceServer_CreateUser.Send] error in sending CreateUser",
		)
		s.logger.Errorw(
			fmt.Sprintf(
				"[UserGrpcServiceServer_CreateUser.Send] id: {%s}, err: %v",
				command.UserId,
				err,
			),
			logger.Fields{"Id": command.UserId},
		)
		return nil, err
	}

	return &userService.CreateUserRes{
		UserId: result.UserID.String(),
	}, nil
}

//
//func (s *UserGrpcServiceServer) UpdateProduct(
//	ctx context.Context,
//	req *userService.UpdateProductReq,
//) (*userService.UpdateProductRes, error) {
//	s.userMetrics.UpdateProductGrpcRequests.Add(ctx, 1, grpcMetricsAttr)
//	span := trace.SpanFromContext(ctx)
//	span.SetAttributes(attribute.Object("Request", req))
//
//	productUUID, err := uuid.FromString(req.GetProductId())
//	if err != nil {
//		badRequestErr := customErrors.NewBadRequestErrorWrap(
//			err,
//			"[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] error in converting uuid",
//		)
//		s.logger.Errorf(
//			fmt.Sprintf(
//				"[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] err: %v",
//				badRequestErr,
//			),
//		)
//		return nil, badRequestErr
//	}
//
//	command, err := updateProductCommandV1.NewUpdateProductWithValidation(
//		productUUID,
//		req.GetName(),
//		req.GetDescription(),
//		req.GetPrice(),
//	)
//	if err != nil {
//		validationErr := customErrors.NewValidationErrorWrap(
//			err,
//			"[ProductGrpcServiceServer_UpdateProduct.StructCtx] command validation failed",
//		)
//		s.logger.Errorf(
//			fmt.Sprintf(
//				"[ProductGrpcServiceServer_UpdateProduct.StructCtx] err: %v",
//				validationErr,
//			),
//		)
//		return nil, validationErr
//	}
//
//	if _, err = mediatr.Send[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](ctx, command); err != nil {
//		err = errors.WithMessage(
//			err,
//			"[ProductGrpcServiceServer_UpdateProduct.Send] error in sending CreateUser",
//		)
//		s.logger.Errorw(
//			fmt.Sprintf(
//				"[ProductGrpcServiceServer_UpdateProduct.Send] id: {%s}, err: %v",
//				command.ProductID,
//				err,
//			),
//			logger.Fields{"Id": command.ProductID},
//		)
//		return nil, err
//	}
//
//	return &userService.UpdateProductRes{}, nil
//}
//
//func (s *UserGrpcServiceServer) GetProductById(
//	ctx context.Context,
//	req *userService.GetProductByIdReq,
//) (*userService.GetProductByIdRes, error) {
//	//// we could use trace manually, but I used grpc middleware for doing this
//	//ctx, span, clean := grpcTracing.StartGrpcServerTracerSpan(ctx, "UserGrpcServiceServer.GetProductById")
//	//defer clean()
//
//	s.userMetrics.GetProductByIdGrpcRequests.Add(ctx, 1, grpcMetricsAttr)
//	span := trace.SpanFromContext(ctx)
//	span.SetAttributes(attribute.Object("Request", req))
//
//	productUUID, err := uuid.FromString(req.GetProductId())
//	if err != nil {
//		badRequestErr := customErrors.NewBadRequestErrorWrap(
//			err,
//			"[ProductGrpcServiceServer_GetProductById.uuid.FromString] error in converting uuid",
//		)
//		s.logger.Errorf(
//			fmt.Sprintf(
//				"[ProductGrpcServiceServer_GetProductById.uuid.FromString] err: %v",
//				badRequestErr,
//			),
//		)
//		return nil, badRequestErr
//	}
//
//	query, err := getProductByIdQueryV1.NewGetProductByIdWithValidation(productUUID)
//	if err != nil {
//		validationErr := customErrors.NewValidationErrorWrap(
//			err,
//			"[ProductGrpcServiceServer_GetProductById.StructCtx] query validation failed",
//		)
//		s.logger.Errorf(
//			fmt.Sprintf(
//				"[ProductGrpcServiceServer_GetProductById.StructCtx] err: %v",
//				validationErr,
//			),
//		)
//		return nil, validationErr
//	}
//
//	queryResult, err := mediatr.Send[*getProductByIdQueryV1.GetProductById, *getProductByIdDtosV1.GetProductByIdResponseDto](
//		ctx,
//		query,
//	)
//	if err != nil {
//		err = errors.WithMessage(
//			err,
//			"[ProductGrpcServiceServer_GetProductById.Send] error in sending GetProductById",
//		)
//		s.logger.Errorw(
//			fmt.Sprintf(
//				"[ProductGrpcServiceServer_GetProductById.Send] id: {%s}, err: %v",
//				query.ProductID,
//				err,
//			),
//			logger.Fields{"Id": query.ProductID},
//		)
//		return nil, err
//	}
//
//	product, err := mapper.Map[*userService.Product](queryResult.Product)
//	if err != nil {
//		err = errors.WithMessage(
//			err,
//			"[ProductGrpcServiceServer_GetProductById.Map] error in mapping product",
//		)
//		return nil, err
//	}
//
//	return &userService.GetProductByIdRes{Product: product}, nil
//}
