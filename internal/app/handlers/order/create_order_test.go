//go:build integration
// +build integration

package order

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	order_service "gitlab.ozon.dev/astoyakin/route256/internal/app/service/order-service"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/currencyStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStatusStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStorage"
	suitePkg "gitlab.ozon.dev/astoyakin/route256/internal/pkg/suite"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
	"reflect"
	"testing"
	"time"
)

type testCaseCreateOrder struct {
	req          *desc.CreateOrderRequest
	expectedResp *desc.CreateOrderResponse
	isError      bool
}

type CreateOrderHeader struct {
	suite.Suite
	db        *suitePkg.DB
	testCases []testCaseCreateOrder
}

func TestCreateOrder(t *testing.T) {
	suite.Run(t, new(CreateOrderHeader))
}

func (c *CreateOrderHeader) SetupSuite() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, "host=localhost port=5432 user=user password=password dbname=order")
	if err != nil {
		c.T().Fatal(err)
	}
	c.db = suitePkg.NewDB(ctx, pool, time.Second*5)

	if err = c.db.TruncateTables(orderTable, orderStatusTable, currencyTable); err != nil {
		c.T().Fatal(err)
	}
	if err = c.db.ImportFixtures("./testdata/currency_table_create.sql", "./testdata/order_table_create.sql"); err != nil {
		c.T().Fatal(err)
	}

	c.testCases = []testCaseCreateOrder{
		{
			req: &desc.CreateOrderRequest{Price: 30,
				Currency:  "USD",
				OrderCode: "f35tfd"},
			expectedResp: &desc.CreateOrderResponse{Id: 1},
			isError:      false,
		},
		{
			req: &desc.CreateOrderRequest{Price: 30,
				Currency:  "RUBBBBB",
				OrderCode: "f35tfd"},
			expectedResp: nil,
			isError:      true,
		},
	}
}

func (c *CreateOrderHeader) TearDownSuite() {
	c.db.Database.GetCloseFunction()
}

func (c *CreateOrderHeader) TestCreateOrder() {
	var (
		ctx             = context.Background()
		repo            = storage.GetRepository(c.db.Database)
		orderRepo       = orderStorage.GetOrderStorage(c.db.Database.GetConnectionPoll())
		orderStatusRepo = orderStatusStorage.GetOrderStatusRepository(c.db.Database.GetConnectionPoll())
		currencyRepo    = currencyStorage.GetCurrencyStorage(c.db.Database.GetConnectionPoll())
	)

	orderService := order_service.NewOrderService(orderRepo, orderStatusRepo, currencyRepo, repo)
	impl := Implementation{orderService: orderService}

	for _, val := range c.testCases {
		resp, err := impl.CreateOrder(ctx, val.req)
		c.Equal(val.isError, err != nil)
		c.Equal(true, reflect.DeepEqual(resp, val.expectedResp))

	}

}
