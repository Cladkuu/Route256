package commandProcessor

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	orderStatusStorage2 "gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStatusStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/config"
	"strconv"
	"strings"
)

const (
	GetAllOrdersCommand = "getAllOrders"
	CreateOrderCommand  = "createOrder"
	CancelOrderCommand  = "cancelOrder"
	GetOrderByIdCommand = "getOrderById"
	ChangeStatusCommand = "changeStatus"
	ResetOrderPrice     = "resetOrderPrice"
	HelpCommand         = "help"
)

type commandFunc func(ctx context.Context, str string) string

type commandProcessor struct {
	commandProcessorMapa map[string]commandFunc
	orderStorage         orderStorage.IOrderStorage
	orderStatusStorage   orderStatusStorage2.IOrderStatusStorage
}

var commandProcessorVariable *commandProcessor

func GetCommandProcessor(OrderStorage orderStorage.IOrderStorage,

	OrderStatusStorage orderStatusStorage2.IOrderStatusStorage) ICommandProcessor {
	if commandProcessorVariable == nil {
		commandProcessorVariable = &commandProcessor{commandProcessorMapa: make(map[string]commandFunc),
			orderStorage:       OrderStorage,
			orderStatusStorage: OrderStatusStorage}

		commandProcessorVariable.commandProcessorMapa[GetAllOrdersCommand] = commandFunc(getAllOrders)
		commandProcessorVariable.commandProcessorMapa[CreateOrderCommand] = commandFunc(createOrder)
		commandProcessorVariable.commandProcessorMapa[CancelOrderCommand] = commandFunc(cancelOrder)
		commandProcessorVariable.commandProcessorMapa[GetOrderByIdCommand] = commandFunc(getOrderById)
		commandProcessorVariable.commandProcessorMapa[ChangeStatusCommand] = commandFunc(changeStatus)
		commandProcessorVariable.commandProcessorMapa[ResetOrderPrice] = commandFunc(resetOrderPrice)
		commandProcessorVariable.commandProcessorMapa[HelpCommand] = commandFunc(help)
	}
	return commandProcessorVariable
}

func (cp *commandProcessor) ProcessCommand(ctx context.Context, command, params string) string {
	if _, ok := cp.commandProcessorMapa[command]; !ok {
		return model.UnknownCommand.Error()
	}

	return cp.commandProcessorMapa[command](ctx, params)
}

func getAllOrders(ctx context.Context, _ string) string {
	var allOrders string
	/*orders, err := commandProcessorVariable.orderStorage.GetAllOrders(ctx)
	if err != nil {
		return err.Error()
	}
	for _, val := range orders {
		allOrders += fmt.Sprintf("GetId: %d, GetStatus: %s, GetPrice: %d, GetCurrency: %s\n", val.GetId(), val.GetStatus(), val.GetPrice(), val.GetCurrency())
	}*/

	return allOrders
}

func createOrder(ctx context.Context, str string) string {
	params := strings.Split(str, " ")

	if len(params) != config.CreateOrderCommandParametersCount {
		return errors.Wrap(model.NotProperCountOfParameters, strconv.Itoa(config.CreateOrderCommandParametersCount)).Error()
	}
	if len(params[1]) != 3 {
		return model.CurrencyError.Error()
	}
	/*id*/ _, err := strconv.Atoi(params[0])
	if err != nil {
		return model.PriceNotIntegerError.Error()
	}
	/*orderId, err := commandProcessorVariable.orderStorage.CreateOrder(ctx, int32(id), params[1], params[3])
	if err != nil {
		return err.Error()
	}
	m := strconv.FormatInt(orderId, 10)*/
	return "order is created. OrderId: " /*+ m*/
}

func getOrderById(ctx context.Context, str string) string {
	if str == "" {
		return errors.Wrap(model.NotProperCountOfParameters, config.GetOrderByIdCommandParametersCount).Error()
	}
	/*id*/ _, err := strconv.Atoi(str)
	if err != nil {
		return model.NotNessesaryTypeOfId.Error()
	}
	/*order, err := commandProcessorVariable.orderStorage.GetOrderById(ctx, int64(id))
	if err != nil {
		return err.Error()
	}*/

	return ""
	/*return fmt.Sprintf("GetId: %d, GetStatus: %s, GetPrice: %d, GetCurrency: %s", order.GetId(), order.GetStatus(), order.GetPrice(), order.GetCurrency())*/

}

func changeStatus(ctx context.Context, str string) string {
	params := strings.Split(str, " ")
	if len(params) != config.ChangeStatusCommandParametersCount {
		return errors.Wrap(model.NotProperCountOfParameters, strconv.Itoa(config.ChangeStatusCommandParametersCount)).Error()
	}
	/*id, err := strconv.Atoi(params[0])
	if err != nil {
		return model.NotNessesaryTypeOfId.Error()
	}
	order, err := commandProcessorVariable.orderStorage.GetOrderById(ctx, int64(id))
	if err != nil {
		return err.Error()
	}
	ok, err := commandProcessorVariable.orderStatusStorage.FindStatuses(order.GetStatus(), params[1])
	if err != nil {
		return err.Error()
	}
	if !ok {
		return model.ForbiddenToPass.Error() + " " + params[1]
	}

	if err = commandProcessorVariable.orderStorage.ChangeStatus(ctx, int64(id), params[1]); err != nil {
		return err.Error()
	}*/

	return "status changed"
}

func resetOrderPrice(ctx context.Context, str string) string {
	if str == "" {
		return errors.Wrap(model.NotProperCountOfParameters, config.ResetOrderPriceCommandParametersCount).Error()
	}
	/*	id, err := strconv.Atoi(str)
		if err != nil {
			return model.NotNessesaryTypeOfId.Error()
		}
		if err = commandProcessorVariable.orderStorage.ResetOrderPrice(ctx, int64(id)); err != nil {
			return err.Error()
		}*/
	return "GetPrice is reset"
}

func cancelOrder(ctx context.Context, str string) string {
	if str == "" {
		return errors.Wrap(model.NotProperCountOfParameters, config.CancelOrderCommandParametersCount).Error()
	}
	/*id, err := strconv.Atoi(str)
	if err != nil {
		return model.NotNessesaryTypeOfId.Error()
	}
	order, err := commandProcessorVariable.orderStorage.GetOrderById(ctx, int64(id))
	if err != nil {
		return err.Error()
	}
	ok, err := commandProcessorVariable.orderStatusStorage.FindStatuses(order.GetStatus(), orderStatusStorage2.CancelledStatus)
	if err != nil {
		return err.Error()
	}
	if !ok {
		return model.ForbiddenToPass.Error() + " " + orderStatusStorage2.CancelledStatus
	}

	if err = commandProcessorVariable.orderStorage.CancelOrder(ctx, int64(id)); err != nil {
		return err.Error()
	}*/

	return "order cancelled"
}

func help(ctx context.Context, _ string) string {
	var helpCommands = "available commands:\n"
	for ind := range commandProcessorVariable.commandProcessorMapa {
		helpCommands += fmt.Sprintf("/%s\n", ind)
	}
	return helpCommands
}
