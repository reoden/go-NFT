package user

import (
	"fmt"

	"github.com/reoden/go-NFT/user/config"
	"github.com/reoden/go-NFT/user/internal/shared/configurations/user/infrastructure"
	"github.com/reoden/go-NFT/user/internal/shared/contracts"
	"github.com/reoden/go-NFT/user/internal/shared/data"
	"github.com/reoden/go-NFT/user/internal/user"

	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/
var UserServiceModule = fx.Module(
	"userfx",
	// Shared Modules
	config.Module,
	infrastructure.Module,
	data.Module,

	// Features Modules
	user.Module,

	// Other provides
	fx.Provide(provideUserMetrics),
)

// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func provideUserMetrics(
	cfg *config.AppOptions,
	meter metric.Meter,
) (*contracts.UserMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	createUserGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_user_grpc_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of create user grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	//updateProductGrpcRequests, err := meter.Float64Counter(
	//	fmt.Sprintf("%s_update_product_grpc_requests_total", cfg.ServiceName),
	//	api.WithDescription("The total number of update product grpc requests"),
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//deleteProductGrpcRequests, err := meter.Float64Counter(
	//	fmt.Sprintf("%s_delete_product_grpc_requests_total", cfg.ServiceName),
	//	api.WithDescription("The total number of delete product grpc requests"),
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//getProductByIdGrpcRequests, err := meter.Float64Counter(
	//	fmt.Sprintf(
	//		"%s_get_product_by_id_grpc_requests_total",
	//		cfg.ServiceName,
	//	),
	//	api.WithDescription(
	//		"The total number of get product by id grpc requests",
	//	),
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//searchProductGrpcRequests, err := meter.Float64Counter(
	//	fmt.Sprintf("%s_search_product_grpc_requests_total", cfg.ServiceName),
	//	api.WithDescription("The total number of search product grpc requests"),
	//)
	//if err != nil {
	//	return nil, err
	//}

	return &contracts.UserMetrics{
		//CreateProductRabbitMQMessages: createProductRabbitMQMessages,
		//GetProductByIdGrpcRequests:    getProductByIdGrpcRequests,
		CreateUserGrpcRequests: createUserGrpcRequests,
		//DeleteProductRabbitMQMessages: deleteProductRabbitMQMessages,
		//DeleteProductGrpcRequests:     deleteProductGrpcRequests,
		//ErrorRabbitMQMessages:         errorRabbitMQMessages,
		//SearchProductGrpcRequests:     searchProductGrpcRequests,
		//SuccessRabbitMQMessages:       successRabbitMQMessages,
		//UpdateProductRabbitMQMessages: updateProductRabbitMQMessages,
		//UpdateProductGrpcRequests:     updateProductGrpcRequests,
	}, nil
}
