package main

import (
	"context"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/consumer/order_event_consumer"
	"gitlab.ozon.dev/astoyakin/route256/internal/config"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/app"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/kafka/consumer"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	app, err := app.GetApp(ctx)
	if err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}
	brokers := strings.Split(os.Getenv("Brokers"), ",")
	/*producerCFG := &kafka.ProducerConfig{
		Brokers: brokers,
		Cfg:     kafka.GetDefaultProducerConfig(),
	}*/
	/*asyncProducer, err := sarama.NewAsyncProducer(producerCFG.Brokers, producerCFG.Cfg)
	if err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}*/

	/*	responseTopic := os.Getenv(config.OrderEventResponseTopic)*/
	consumerTopic := os.Getenv(config.OrderEventTopic)
	logger.GlobalLogger.Info("Topic Name: " + consumerTopic)

	/*kafkaProducer := kafka.NewKafkaProducer(asyncProducer,
	producerCFG,
	map[string]string{
		"responseTopic": responseTopic,
		"ourTopic":      consumerTopic,
	})*/
	kafkaProducer, err := app.InitKafkaProducer()
	if err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}

	defer kafkaProducer.Close()

	handler := order_event_consumer.New(app.GetOrderService(), kafkaProducer, 3) // TODO add from env

	groupId := os.Getenv(config.KafkaGroupId)
	logger.GlobalLogger.Info("groupId: " + groupId)
	// TODO добавить эти данные в env

	consumerGroup, err := consumer.NewConsumerGroup(ctx,
		&consumer.Config{
			Topics:       []string{consumerTopic},
			Brokers:      brokers,
			GroupId:      groupId,
			Handler:      handler,
			SaramaConfig: consumer.GetDefaultConsumerGroupConfig(),
		})
	if err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}
	defer consumerGroup.Close()

	if err = consumerGroup.Run(ctx); err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}
}
