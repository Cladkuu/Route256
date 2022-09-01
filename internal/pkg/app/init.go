package app

import (
	"context"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/handlers/order"
	"gitlab.ozon.dev/astoyakin/route256/internal/bot"
	configPkg "gitlab.ozon.dev/astoyakin/route256/internal/config"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/closer"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/db"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/middleware"
	order_api "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
	pb "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"strings"
	"time"
)

func (a *App) initApp(ctx context.Context) error {
	funcs := []func(context.Context) error{
		a.initEnv,
		a.initCloser,
		a.initDB,
		a.initTracer,
		a.initTelegramBot,
		a.initGrpcServer,
		a.initHTTPServer,
	}
	for _, val := range funcs {
		if err := val(ctx); err != nil {
			return err
		}

	}

	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	if a.telegramBot == nil {
		bot, err := bot.NewTelegramBot(a.initCommandProcessor())
		if err != nil {
			return err
		}
		a.telegramBot = bot
	}
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {

	orderApiImpl := order.NewOrderApiImplementation(a.initOrderService())
	a.orderApiImpl = orderApiImpl

	opts := []grpc.ServerOption{
		/*grpc.UnaryInterceptor(middleware.NewGrpcInterceptor()),*/
		grpc.ChainUnaryInterceptor(middleware.NewGrpcInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(a.tracer))),
	}

	grpcServer := grpc.NewServer(opts...)
	order_api.RegisterOrderApiServer(grpcServer, orderApiImpl)

	a.grpcServers = grpcServer
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	/*ctx, cancel := context.WithCancel(ctx)
	defer cancel()*/

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterOrderApiHandlerFromEndpoint(ctx, mux, os.Getenv(configPkg.GrpcServerAddress), opts); err != nil {
		return err
	}
	/*if err := pb.RegisterOrderApiHandlerServer(ctx, mux, a.orderApiImpl); err != nil {
		return err
	}*/

	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)
	fs := http.FileServer(http.Dir(os.Getenv(configPkg.SwaggerComponentsFilepath)))
	httpMux.Handle(os.Getenv(configPkg.SwaggerUiEndpoint), http.StripPrefix(os.Getenv(configPkg.SwaggerUiEndpoint), fs))

	a.httpServers = mux
	a.SwaggerUiServer = httpMux
	return nil
}

func headerMatcherREST(header string) (string, bool) {
	header = strings.ToLower(header)
	switch header {
	case middleware.UserIdHeader:
		return header, true
	default:
		return header, false
	}

}

func (a *App) initDB(ctx context.Context) error {
	pool, err := pgxpool.Connect(ctx, os.Getenv(configPkg.DbConnString)) // Add conn string dbname=order host=localhost port=5432 user=user password=password dbname=order sslmode=disable
	if err != nil {
		return err
	}

	if err = pool.Ping(ctx); err != nil {
		return err
	}

	/*	config := pool.Config()
		config.MaxConnIdleTime = 2
		config.MaxConns = 20
		config.MinConns = 2
		config.MaxConnLifetime = time.Hour */

	a.psql = db.NewDB(ctx, pool, time.Duration(configPkg.ContextCancellationTimeout)*time.Millisecond)
	a.closer.Add(a.psql.GetCloseFunction)

	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	a.closer = closer.NewCloser()
	a.closer.Add(logger.GlobalLogger.Sync) // close logger
	return nil
}

func (a *App) initEnv(ctx context.Context) error {
	if err := godotenv.Load("config.env"); err != nil {
		return err
	}

	return nil
}
