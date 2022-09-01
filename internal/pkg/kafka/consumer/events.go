package consumer

const (
	KafkaCreateOrderEvent       = "CreateOrder"
	KafkaResetOrderPriceEvent   = "ResetOrderPrice"
	KafkaCancelOrderEvent       = "CancelOrder"
	KafkaChangeOrderStatusEvent = "ChangeOrder"
)

/*var (
	KafkaEventMap = map[string]struct{}{
		KafkaCreateOrderEvent:       {},
		KafkaResetOrderPriceEvent:   {},
		KafkaChangeOrderStatusEvent: {},
		KafkaCancelOrderEvent:       {},
	}

)*/
