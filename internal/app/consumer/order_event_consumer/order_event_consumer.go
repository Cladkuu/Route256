package order_event_consumer

import (
	"context"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	order_service "gitlab.ozon.dev/astoyakin/route256/internal/app/service/order-service"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/kafka"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/kafka/consumer"
)

type orderEventConsumer struct {
	orderService order_service.IOrder
	kafkaProducer/**kafka.KafkaProducer*/ kafka.IKafkaProducer
	maxRetryCount int8
	// TODO add kafka producer
}

func New(OrderService order_service.IOrder,
	KafkaProducer kafka.IKafkaProducer,
	MaxRetryCount int8) *orderEventConsumer {
	return &orderEventConsumer{
		orderService:  OrderService,
		kafkaProducer: KafkaProducer,
		maxRetryCount: MaxRetryCount,
	}
}

func (con *orderEventConsumer) Handle(ctx context.Context, message *sarama.ConsumerMessage) error {
	if message == nil {
		return nil
	}

	mapa := make(map[string]jsoniter.RawMessage)

	if err := jsoniter.Unmarshal(message.Value, &mapa); err != nil {
		// TODO add logging and return nil
		return err
	}
	event := consumer.CommonOrderEvent{}
	if err := jsoniter.Unmarshal(mapa["msg"], &event); err != nil {
		// TODO add logging and return nil
		return err
	}

	// Проверка, что кол-во доступных повторных попыток не превышено
	if con.maxRetryCount <= event.RetryCount {
		if err := con.kafkaProducer.SendMessageToTopic(ctx, con.kafkaProducer.GetTopic(event.Event), kafka.KafkaEventResponse{
			Event:  event.Event,
			Status: kafka.FailedResponse,
			TaskId: event.TaskId,
		}); err != nil {
			return err
		} // TODO add real topic and msg
		return nil
	}

	switch event.Event {
	case consumer.KafkaCancelOrderEvent:
		cancel := consumer.CancelOrderEvent{}
		if err := jsoniter.Unmarshal(message.Value, &cancel); err != nil {
			// TODO add logging
			return err
		}

		if err := con.orderService.CancelOrder(ctx, cancel.Id); err != nil {
			cancel.Msg.RetryCount++
			if err = con.kafkaProducer.SendMessageToTopic(ctx, "order_event", cancel); err != nil { // TODO add real topic
				return err
			}
			return nil
		}

		if err := con.kafkaProducer.SendMessageToTopic(ctx, con.kafkaProducer.GetTopic(event.Event), kafka.KafkaEventResponse{
			Status: kafka.SuccessResponse,
			Event:  kafka.KafkaCancelOrderResponseEvent,
			TaskId: event.TaskId,
		}); err != nil { // TODO add real topic and msg
			return err
		}

	case consumer.KafkaResetOrderPriceEvent:
		reset := consumer.ResetOrderPriceEvent{}
		if err := jsoniter.Unmarshal(message.Value, &reset); err != nil {
			// TODO add logging
			return err
		}

		if err := con.orderService.ResetOrderPrice(ctx, reset.Id); err != nil {
			reset.Msg.RetryCount++
			if err = con.kafkaProducer.SendMessageToTopic(ctx, "order_event", reset); err != nil { // TODO add real topic
				return err
			}
			return nil
		}

		if err := con.kafkaProducer.SendMessageToTopic(ctx, con.kafkaProducer.GetTopic(event.Event), kafka.KafkaEventResponse{
			Status: kafka.SuccessResponse,
			Event:  kafka.KafkaResetOrderPriceResponseEvent,
			TaskId: event.TaskId,
		}); err != nil { // TODO add real topic and msg
			return err
		}

	case consumer.KafkaChangeOrderStatusEvent:
		change := consumer.ChangeOrderStatus{}
		if err := jsoniter.Unmarshal(message.Value, &change); err != nil {
			// TODO add logging
			return err
		}

		if err := con.orderService.ChangeOrderStatus(ctx, change.Id, change.Status); err != nil {
			change.Msg.RetryCount++
			if err = con.kafkaProducer.SendMessageToTopic(ctx, "order_event", change); err != nil { // TODO add real topic
				return err
			}
			return nil
		}

		if err := con.kafkaProducer.SendMessageToTopic(ctx, con.kafkaProducer.GetTopic(event.Event), kafka.KafkaEventResponse{
			Status: kafka.SuccessResponse,
			Event:  kafka.KafkaChangeOrderStatusResponseEvent,
			TaskId: event.TaskId,
		}); err != nil { // TODO add real topic and msg
			return err
		}

	case consumer.KafkaCreateOrderEvent:
		create := consumer.OrderCreateEvent{}
		if err := jsoniter.Unmarshal(message.Value, &create); err != nil {
			// TODO add logging
			return err
		}

		id, err := con.orderService.CreateOrder(ctx, create.Price, create.Currency, create.OrderCode)
		if err != nil {
			create.Msg.RetryCount++
			if err = con.kafkaProducer.SendMessageToTopic(ctx, "order_event", create); err != nil { // TODO add real topic
				return err
			}
			return nil
		}

		if err = con.kafkaProducer.SendMessageToTopic(ctx /*con.kafkaProducer.GetTopic(event.Event)*/, "order_event_response", kafka.CreateOrderSuccessResponse{Common: kafka.KafkaEventResponse{
			Status: kafka.SuccessResponse,
			Event:  kafka.KafkaCreateOrderResponseEvent,
			TaskId: event.TaskId,
		},
			Id: id}); err != nil { // TODO add real topic and msg
			return err
		}

	}

	return nil
}
