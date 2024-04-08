package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/project-alvarium/alvarium-sdk-go/pkg"
	SdkConfig "github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/factories"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/ones-demo-2024/internal/bootstrap"
	"github.com/project-alvarium/ones-demo-2024/internal/config"
	"github.com/project-alvarium/ones-demo-2024/internal/transitor"
	"log/slog"
	"os"
)

func main() {
	// Load config
	var configPath string
	flag.StringVar(&configPath,
		"cfg",
		"./res/config.json",
		"Path to JSON configuration file.")
	flag.Parse()

	fileFormat := config.GetFileExtension(configPath)
	reader, err := config.NewReader(fileFormat)
	if err != nil {
		tmpLog := factories.NewLogger(SdkConfig.LoggingInfo{MinLogLevel: slog.LevelError})
		tmpLog.Error(err.Error())
		os.Exit(1)
	}

	cfg := config.ApplicationConfig{}
	err = reader.Read(configPath, &cfg)
	if err != nil {
		tmpLog := factories.NewLogger(SdkConfig.LoggingInfo{MinLogLevel: slog.LevelError})
		tmpLog.Error(err.Error())
		os.Exit(1)
	}

	logger := factories.NewLogger(cfg.Logging)
	logger.Write(slog.LevelDebug, "config loaded successfully")
	logger.Write(slog.LevelDebug, cfg.AsString())

	// List of annotators driven from config, eventually support dist. policy.
	var annotators []interfaces.Annotator
	for _, t := range cfg.Sdk.Annotators {
		instance, err := factories.NewAnnotator(t, cfg.Sdk)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		annotators = append(annotators, instance)
	}
	sdk := pkg.NewSdk(annotators, cfg.Sdk, logger)

	r := mux.NewRouter()
	transitor.LoadRestRoutes(r, sdk, logger)
	ctx, cancel := context.WithCancel(context.Background())
	bootstrap.Run(
		ctx,
		cancel,
		cfg,
		[]bootstrap.BootstrapHandler{
			transitor.NewHttpServer(r, cfg.Endpoint, logger).BootstrapHandler,
			sdk.BootstrapHandler,
		})
}
