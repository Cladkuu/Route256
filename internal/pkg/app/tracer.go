package app

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	/*"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"*/)

/*var (
	globalTracer opentracing.Tracer
	tracerCloser io.Closer
)*/

func (a *App) initTracer(ctx context.Context) error {
	cfg := config.Configuration{
		ServiceName: a.name,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{LogSpans: true},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return err
	}
	a.tracer = tracer
	a.closer.Add(closer.Close)
	opentracing.SetGlobalTracer(tracer)

	return nil
}
