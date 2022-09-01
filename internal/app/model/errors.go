package model

import "errors"

var (
	UnknownStatus              = errors.New("unknown GetStatus")
	UnknownCommand             = errors.New("unknown Command")
	NotNessesaryTypeOfId       = errors.New("identifier type must be integer")
	NotProperCountOfParameters = errors.New(" parameter(s) must be passed\nCommand didn't processed")
	NotFoundError              = errors.New("Not Found")
	ForbiddenToPass            = errors.New("Forbidden to pass over in new Status")
	CurrencyError              = errors.New("Error. Second parameter is not a Currency")
	PriceNotIntegerError       = errors.New("GetPrice should be integer")
	PriceIsLessZeroErr         = errors.New("GetPrice is less 0")
	UnknownCurrency            = errors.New("unknown currency")
	NotUniqueOrderCode         = errors.New("Not Unique Order Code")
)

// TODO убрать не нужные ошибки
// TODO сделать мапу ошибок
