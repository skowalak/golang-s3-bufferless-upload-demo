package main

import (
	"context"
	"file-upload-demo/aws"
	"file-upload-demo/config"
	"file-upload-demo/httphandler"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	// logger sync is called in app shutdown hook
	slogger := logger.Sugar()

	return slogger
}

func main() {
	fx.New(
		fx.Provide(config.ProvideConfig),
		fx.Provide(ProvideLogger),
		fx.Provide(http.NewServeMux),
		fx.Provide(aws.NewS3),
		fx.Invoke(httphandler.New),
		fx.Invoke(runApp),
	).Run()
}

func runApp(lifecycle fx.Lifecycle, logger *zap.SugaredLogger, cfg *config.Config, mux *http.ServeMux) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go http.ListenAndServe(cfg.ApplicationConfig.Address, mux)
				return nil
			},
			OnStop: func(context.Context) error {
				return logger.Sync()
			},
		},
	)
}
