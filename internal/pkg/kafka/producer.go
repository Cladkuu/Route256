package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
	"time"
)

type KafkaProducer struct {
	AsyncProducer sarama.AsyncProducer
	config        *ProducerConfig
	kafkaTopics   map[string]string
}

func NewKafkaProducer(asyncProducer sarama.AsyncProducer,
	Config *ProducerConfig,
	KafkaTopics map[string]string) IKafkaProducer {

	go func() {

		for err := range asyncProducer.Errors() {
			logger.GlobalLogger.Error("message doesnt send: " + err.Error())
			time.Sleep(200 * time.Millisecond)
			asyncProducer.Input() <- &sarama.ProducerMessage{
				Topic: err.Msg.Topic,
				Value: err.Msg.Value,
			}
		}

	}()

	return &KafkaProducer{
		AsyncProducer: asyncProducer,
		config:        Config,
		kafkaTopics:   KafkaTopics,
	}
}

func (kp *KafkaProducer) GetTopic(key string) string {
	return kp.kafkaTopics[key]
}

func (k *KafkaProducer) Close() error {
	k.AsyncProducer.Close()
	return nil
}

func (k *KafkaProducer) SendMessageToTopic(ctx context.Context, topic string, msg interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendMessageToTopic")
	defer span.Finish()
	span.LogFields(log.String("topic", topic))

	byteArray, err := jsoniter.Marshal(msg)
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		return err
	}

	k.AsyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(byteArray),
	}
	return nil
}
