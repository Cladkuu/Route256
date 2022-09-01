package kafka

import (
	"context"
)

type IKafkaProducer interface {
	GetTopic(key string) string
	Close() error
	SendMessageToTopic(ctx context.Context, topic string, msg interface{}) error
}
