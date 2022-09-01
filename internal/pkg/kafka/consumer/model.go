package consumer

type OrderCreateEvent struct {
	Msg        CommonOrderEvent `json:"msg"`
	Price      int32            `json:"price"`
	Currency   string           `json:"currency"`
	OrderCode  string           `json:"orderCode"`
	RetryCount int8             `json:"retryCount"`
}

type CancelOrderEvent struct {
	Msg CommonOrderEvent `json:"msg"`
	Id  int64            `json:"id"`
}

type ResetOrderPriceEvent struct {
	Msg CommonOrderEvent `json:"msg"`
	Id  int64            `json:"id"`
}

type ChangeOrderStatus struct {
	Msg    CommonOrderEvent `json:"msg"`
	Id     int64            `json:"id"`
	Status string           `json:"status"`
}

type CommonOrderEvent struct {
	TaskId     int64  `json:"taskId"`
	Event      string `json:"event"`
	RetryCount int8   `json:"retryCount"`
}
