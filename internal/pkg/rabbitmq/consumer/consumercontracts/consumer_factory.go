package consumercontracts

import (
	"github.com/reoden/go-echo-template/pkg/core/messaging/consumer"
	messagingTypes "github.com/reoden/go-echo-template/pkg/core/messaging/types"
	"github.com/reoden/go-echo-template/pkg/rabbitmq/consumer/configurations"
	"github.com/reoden/go-echo-template/pkg/rabbitmq/types"
)

type ConsumerFactory interface {
	CreateConsumer(
		consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
		isConsumedNotifications ...func(message messagingTypes.IMessage),
	) (consumer.Consumer, error)

	Connection() types.IConnection
}
