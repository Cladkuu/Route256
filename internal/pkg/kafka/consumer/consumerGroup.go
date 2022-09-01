package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
	"time"
)

type consumerGroup struct {
	consumer      *Consumer
	config        *Config
	consumerGroup sarama.ConsumerGroup
}

func NewConsumerGroup(_ context.Context, cfg *Config) (*consumerGroup, error) {
	client, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupId, cfg.SaramaConfig)
	if err != nil {
		return nil, err
	}

	csGr := consumerGroup{
		config:        cfg,
		consumerGroup: client,
	}

	if cfg.Handler != nil {
		csGr.consumer = NewConsumer(cfg.Handler)
	} // TODO add return error while handler==nil

	return &csGr, nil
}

func (cs *consumerGroup) Run(ctx context.Context) error {
	logger.GlobalLogger.Info("start consumer run...")

	for {
		if err := cs.consumerGroup.Consume(ctx, cs.config.Topics, cs.consumer); err != nil {
			time.Sleep(time.Millisecond * 150)
			logger.GlobalLogger.Info("consume error: " + err.Error())
			if err == sarama.ErrClosedConsumerGroup {
				return err
			}
		}
	}
}

func (cs *consumerGroup) Close() error {
	logger.GlobalLogger.Info("close consumer Group")

	cs.consumerGroup.Close()
	return nil
}
