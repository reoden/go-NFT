package bus

import (
	consumer2 "github.com/reoden/go-echo-template/pkg/core/messaging/consumer"
	"github.com/reoden/go-echo-template/pkg/core/messaging/producer"
)

type Bus interface {
	producer.Producer
	consumer2.BusControl
	consumer2.ConsumerConnector
}
