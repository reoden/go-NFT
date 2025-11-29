package configurations

import (
	consumerConfigurations "github.com/reoden/go-NFT/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/reoden/go-NFT/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfiguration struct {
	ProducersConfigurations []*producerConfigurations.RabbitMQProducerConfiguration
	ConsumersConfigurations []*consumerConfigurations.RabbitMQConsumerConfiguration
}
