package configurations

import (
	consumerConfigurations "github.com/reoden/go-echo-template/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/reoden/go-echo-template/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfiguration struct {
	ProducersConfigurations []*producerConfigurations.RabbitMQProducerConfiguration
	ConsumersConfigurations []*consumerConfigurations.RabbitMQConsumerConfiguration
}
