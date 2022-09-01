package app

import order_service "gitlab.ozon.dev/astoyakin/route256/internal/app/service/order-service"

func (a *App) GetOrderService() order_service.IOrder {
	if a.orderService == nil {
		order_service.NewOrderService(a.initOrderStorage(),
			a.initOrderStatusStorage(),
			a.initCurrencyStorage(),
			a.initCommonRepository())
	}
	return a.orderService
}
