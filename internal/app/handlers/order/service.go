package order

import (
	orderServicePkg "gitlab.ozon.dev/astoyakin/route256/internal/app/service/order-service"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

type Implementation struct {
	desc.UnimplementedOrderApiServer
	orderService orderServicePkg.IOrder
}

func NewOrderApiImplementation(OrderService orderServicePkg.IOrder) *Implementation {
	return &Implementation{
		orderService: OrderService,
	}
}
