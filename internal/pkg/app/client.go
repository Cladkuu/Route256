package app

import (
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/astoyakin/route256/internal/config"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/kafka"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/kafka/consumer"
	"os"
	"strings"
)

func (a *App) InitKafkaProducer() (kafka.IKafkaProducer, error) {
	if a.kafkaProducer == nil {
		brokers := strings.Split(os.Getenv( /*"KafkaBrokers"*/ config.Brokers), ",")
		cfg := kafka.ProducerConfig{Brokers: brokers,
			Cfg: kafka.GetDefaultProducerConfig()}
		producer, err := sarama.NewAsyncProducer(brokers, cfg.Cfg)
		if err != nil {
			return nil, err
		}
		topic := os.Getenv(config.OrderEventResponseTopic)

		a.kafkaProducer = kafka.NewKafkaProducer(producer, &cfg, map[string]string{
			consumer.KafkaCreateOrderEvent:       topic,
			consumer.KafkaCancelOrderEvent:       topic,
			consumer.KafkaChangeOrderStatusEvent: topic,
			consumer.KafkaResetOrderPriceEvent:   topic,
		}) // TODO подумать над тем, чтобы эвенты вынести в env в виде массива
	}
	a.closer.Add(a.kafkaProducer.Close)
	return a.kafkaProducer, nil
}
