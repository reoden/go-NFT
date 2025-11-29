package rabbitmq

import (
	"github.com/reoden/go-NFT/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/reoden/go-NFT/pkg/rabbitmq/producer/configurations"
	"github.com/reoden/go-NFT/user/internal/user/features/creatingproduct/v1/events/integrationevents"
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
