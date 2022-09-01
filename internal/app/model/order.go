package model

var orderId int64

type Order struct {
	Id        int64  `db:"id"`
	Status    string `db:"status"`
	Price     int32  `db:"price"`
	Currency  string `db:"currency"`
	orderCode string `db:"-"`
}

func CreateOrder(Price int32, Currency, status string, OrderCode string) *Order {
	orderId++
	return &Order{Id: orderId,
		Status:    status,
		Price:     Price,
		Currency:  Currency,
		orderCode: OrderCode}
}

func (o *Order) GetId() int64 {
	return o.Id
}

func (o *Order) GetStatus() string {
	return o.Status
}

func (o *Order) SetStatus(status string) {
	o.Status = status
}

func (o *Order) GetPrice() int32 {
	return o.Price
}

func (o *Order) ResetPrice() {
	o.Price = 0
}

func (o *Order) GetCurrency() string {
	return o.Currency
}

func (o *Order) SetPrice(Price int32) {
	o.Price = Price
}
