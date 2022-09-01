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

const (
	orderTable         = "public.orders"
	orderStatusTable   = "public.status"
	statusMappingTable = "public.status_mapping"
	currencyTable      = "public.currency"
)

type testCaseCancelOrder struct {
	req          *desc.CancelOrderRequest
	expectedResp *desc.CancelOrderResponse
	isError      bool
}

type CancelOrderHeader struct {
	suite.Suite
	db        *suitePkg.DB
	testCases []testCaseCancelOrder
}

func TestOrderCancellation(t *testing.T) {
	suite.Run(t, new(CancelOrderHeader))
}

func (c *CancelOrderHeader) SetupSuite() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, "host=localhost port=6432 user=user password=password dbname=order")
	if err != nil {
		c.Fail("get pool connection %s", err.Error())
	}
	c.db = suitePkg.NewDB(ctx, pool, time.Second)

	if err = c.db.TruncateTables(orderTable, statusMappingTable, orderStatusTable); err != nil {
		c.T().Fatal(err)
	}
	if err = c.db.ImportFixtures("./testdata/order_status_table_create.sql", "./testdata/order_table_create.sql", "./testdata/status_mapping_table_create.sql", "./testdata/order_table_insert.sql"); err != nil {
		c.Fail("can't import fixtures", err.Error())
	}

	c.testCases = []testCaseCancelOrder{
		{
			req:          &desc.CancelOrderRequest{Id: 1},
			expectedResp: &desc.CancelOrderResponse{Result: true},
			isError:      false,
		},
		{
			req:          &desc.CancelOrderRequest{Id: 1},
			expectedResp: nil,
			isError:      true,
		},
	}
}

func (c *CancelOrderHeader) TearDownSuite() {
	c.db.Database.GetCloseFunction()
}

func (c *CancelOrderHeader) TestCancelOrder() {
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
		resp, err := impl.CancelOrder(ctx, val.req)
		c.Equal(val.isError, err != nil)
		c.Equal(true, reflect.DeepEqual(resp, val.expectedResp))

	}

}
