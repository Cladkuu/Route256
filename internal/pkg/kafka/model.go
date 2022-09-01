package kafka

type KafkaEventResponse struct {
	Event  string `json:"event"`
	Status string `json:"status"`
	TaskId int64  `json:"taskId"`
}

type CreateOrderSuccessResponse struct {
	Common KafkaEventResponse `json:"common"`
	Id     int64              `json:"orderId"`
}

const (
	SuccessResponse = "SUCCESS"
	FailedResponse  = "FAILED"
) // TODO вывести данные константы в тип
