//go:generate mockery --name=(.+)Mock --case=underscore

package storage

import (
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/currencyStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStatusStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStorage"
)

type CurrencyRepositoryMock interface {
	currencyStorage.ICurrency
}

type OrderStatusRepositoryMock interface {
	orderStatusStorage.IOrderStatusStorage
}

type OrderRepositoryMock interface {
	orderStorage.IOrderStorage
}

type RepositoryMock interface {
	IRepository
}
