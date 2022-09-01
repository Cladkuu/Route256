package convertors

import (
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func ConvertToAllOrders(orders []*model.Order) *desc.GetAllOrdersResponse {
	result := &desc.GetAllOrdersResponse{Order: make([]*desc.Order, 0, len(orders))}
	for _, val := range orders {
		result.Order = append(result.Order, &desc.Order{
			Status:   val.GetStatus(),
			Price:    val.GetPrice(),
			Currency: val.GetCurrency(),
		})
	}
	return result
}
