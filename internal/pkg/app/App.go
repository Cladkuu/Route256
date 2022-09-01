package app

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/handlers/order"
	order_service "gitlab.ozon.dev/astoyakin/route256/internal/app/service/order-service"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/currencyStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStatusStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/bot"
	"gitlab.ozon.dev/astoyakin/route256/internal/bot/commandProcessor"
	"gitlab.ozon.dev/astoyakin/route256/internal/config"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/closer"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/db"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/kafka"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

var (
	application *App
)

type App struct {
	name string
	// storage
	orderStorage       orderStorage.IOrderStorage
	currencyStorage    currencyStorage.ICurrency
	orderStatusStorage orderStatusStorage.IOrderStatusStorage
	commonRepository   storage.IRepository

	// bot
	orderService     order_service.IOrder
	commandProcessor commandProcessor.ICommandProcessor

	orderApiImpl *order.Implementation
	// servers
	grpcServers     *grpc.Server
	httpServers     *runtime.ServeMux
	SwaggerUiServer *http.ServeMux
	telegramBot     *bot.TelegramBot

	// db
	psql db.IDB

	// closer
	closer closer.ICloser

	// tracer
	tracer opentracing.Tracer

	// kafka producer
	kafkaProducer kafka.IKafkaProducer
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{name: config.AppName}
	if err := app.initApp(ctx); err != nil {
		return nil, err
	}
	application = app
	return app, nil
}

func (a *App) Run(_ context.Context) error {
	go a.telegramBot.StartWorking()

	listener, err := net.Listen(os.Getenv(config.GrpcNetwork) /*"localhost:8081"*/, os.Getenv(config.GrpcServerAddress))
	if err != nil {
		return err
	}
	errGr, _ := errgroup.WithContext(context.Background())

	errGr.Go(func() error {
		if err = a.grpcServers.Serve(listener); err != nil {
			return err
		}
		return nil
	})

	errGr.Go(func() error {

		if err = http.ListenAndServe(os.Getenv(config.HttpServerAddress), a.httpServers); err != nil {
			return err
		}
		return nil
	})

	errGr.Go(func() error {
		if err = http.ListenAndServe(os.Getenv(config.SwaggerUiAddress), a.SwaggerUiServer); err != nil {
			return err
		}
		return nil
	})

	if err = errGr.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) GetCloser() closer.ICloser {
	if a.closer == nil {
		a.closer = closer.NewCloser()
	}
	return a.closer
}

func GetApp(ctx context.Context) (*App, error) {
	if application == nil {
		app, err := NewApp(ctx)
		if err != nil {
			return nil, err
		}
		application = app
	}
	return application, nil
}
