package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
)

type IConsumer interface {
	Handle(ctx context.Context, message *sarama.ConsumerMessage) error
}

type Consumer struct {
	handler IConsumer
}

func NewConsumer(Handler IConsumer) *Consumer {
	return &Consumer{
		handler: Handler,
	}
}

// Setup ...
func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup ...
func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			logger.GlobalLogger.Error("context is exceeded")
			return nil
		case msg := <-claim.Messages():
			ctx := context.Background()
			logger.GlobalLogger.Info("msg")
			span, spanCtx := opentracing.StartSpanFromContext(ctx, "kafka consumer")
			if err := c.handler.Handle(spanCtx, msg); err != nil {
				logger.GlobalLogger.Error(err.Error())
				// TODO добавить обработку сообщения
			}

			span.Finish()
			session.MarkMessage(msg, "")

			// TODO Add Message Handler
		}

	}
}
