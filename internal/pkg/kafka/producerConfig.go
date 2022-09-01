package kafka

import (
	"github.com/Shopify/sarama"
	"time"
)

type ProducerConfig struct {
	Brokers []string
	Cfg     *sarama.Config
}

func GetDefaultProducerConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Producer.Timeout = time.Second // TODO add variable from config
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 3 // TODO add variable from config

	return cfg
}
