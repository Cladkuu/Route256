package middleware

import (
	"context"
	"database/sql"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	UserIdHeader   = "user-id"
	EmptyUserIdErr = "UserId must be passed"
)

var (
	clientsChan = make(chan struct{}, 10)
)

func NewGrpcInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		clientsChan <- struct{}{}
		defer func() {
			<-clientsChan
		}()
		/*md*/ _, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, EmptyUserIdErr)
		}
		/*if md.Get(UserIdHeader) == nil || md.Get(UserIdHeader)[0] == "" {
			return nil, status.Error(codes.PermissionDenied, EmptyUserIdErr)
		}*/

		// Создание трасировки запроса
		span, ctxSpan := opentracing.StartSpanFromContext(ctx, info.FullMethod)
		defer span.Finish()
		span.LogFields(log.Object("request body", req))
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			logger.GlobalLogger.Info("traceId: " + sc.SpanID().String())
		}

		response, err := handler(ctxSpan, req)

		if err != nil {
			err = handeError(err)
			span.LogFields(log.Object("error", err))
			return nil, err
		}
		span.LogFields(log.Object("response body", response))
		return response, err

	}
}

func handeError(err error) error {
	if errors.Is(err, model.NotFoundError) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, model.UnknownStatus) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, model.ForbiddenToPass) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, err.Error())
	} else if errors.Is(err, model.PriceIsLessZeroErr) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())

}
