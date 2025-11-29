package producercontracts

import (
	"github.com/reoden/go-NFT/pkg/core/messaging/producer"
	types2 "github.com/reoden/go-NFT/pkg/core/messaging/types"
	"github.com/reoden/go-NFT/pkg/rabbitmq/producer/configurations"
)

type ProducerFactory interface {
	CreateProducer(
		rabbitmqProducersConfiguration map[string]*configurations.RabbitMQProducerConfiguration,
		isProducedNotifications ...func(message types2.IMessage),
	) (producer.Producer, error)
}
