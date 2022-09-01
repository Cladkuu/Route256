package app

import (
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage"
	currencyStorage2 "gitlab.ozon.dev/astoyakin/route256/internal/app/storage/currencyStorage"
	orderStatusStorage2 "gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStatusStorage"
	orderStorage2 "gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStorage"
)

func (a *App) initOrderStorage() orderStorage2.IOrderStorage {
	if a.orderStorage == nil {
		a.orderStorage = orderStorage2.GetOrderStorage(a.psql.GetConnectionPoll())
	}
	return a.orderStorage
}

func (a *App) initOrderStatusStorage() orderStatusStorage2.IOrderStatusStorage {
	if a.orderStatusStorage == nil {
		a.orderStatusStorage = orderStatusStorage2.GetOrderStatusRepository(a.psql.GetConnectionPoll())
	}
	return a.orderStatusStorage
}

func (a *App) initCurrencyStorage() currencyStorage2.ICurrency {
	if a.currencyStorage == nil {
		a.currencyStorage = currencyStorage2.GetCurrencyStorage(a.psql.GetConnectionPoll())
	}
	return a.currencyStorage
}

func (a *App) initCommonRepository() storage.IRepository {
	if a.commonRepository == nil {
		a.commonRepository = storage.GetRepository(a.psql)
	}
	return a.commonRepository
}
