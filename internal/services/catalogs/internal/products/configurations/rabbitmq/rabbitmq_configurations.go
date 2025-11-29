package rabbitmq

import (
	"github.com/reoden/go-NFT/catalogs/internal/products/features/creatingproduct/v1/events/integrationevents"
	"github.com/reoden/go-NFT/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/reoden/go-NFT/pkg/rabbitmq/producer/configurations"
)

func ConfigProductsRabbitMQ(
	builder configurations.RabbitMQConfigurationBuilder,
) {
	builder.AddProducer(
		integrationevents.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		},
	)
}
