package main

import (
	"context"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/app"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
)

func main() {
	ctx := context.Background()
	App, err := app.NewApp(ctx)
	if err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}
	closer := App.GetCloser()
	defer closer.CloseAll()

	if err = App.Run(ctx); err != nil {
		logger.GlobalLogger.Fatal(err.Error())
	}
}
