package app

import (
	order_service "gitlab.ozon.dev/astoyakin/route256/internal/app/service/order-service"
	"gitlab.ozon.dev/astoyakin/route256/internal/bot/commandProcessor"
)

func (a *App) initCommandProcessor() commandProcessor.ICommandProcessor {
	if a.commandProcessor == nil {
		a.commandProcessor = commandProcessor.GetCommandProcessor(a.initOrderStorage(),
			a.initOrderStatusStorage())
	}
	return a.commandProcessor
}

func (a *App) initOrderService() order_service.IOrder {
	if a.orderService == nil {
		a.orderService = order_service.NewOrderService(a.initOrderStorage(),
			a.initOrderStatusStorage(),
			a.initCurrencyStorage(),
			a.initCommonRepository())
	}
	return a.orderService
}
