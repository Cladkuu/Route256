package consumer

import (
	"github.com/Shopify/sarama"
	"time"
)

type Config struct {
	Brokers      []string
	GroupId      string
	Topics       []string
	SaramaConfig *sarama.Config
	Handler      IConsumer
}

func GetDefaultConsumerGroupConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V1_1_1_0
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.MaxWaitTime = 300 * time.Millisecond
	cfg.Consumer.MaxProcessingTime = 3 * time.Second

	return cfg
}
