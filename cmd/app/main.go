package main

import (
	"context"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/chaindead/go-template/internal/app"
	"github.com/chaindead/go-template/internal/config"
	"github.com/chaindead/go-template/internal/logger"
	metricsrv "github.com/chaindead/go-template/internal/metrics/server"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := config.Parse(); err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := logger.Setup(); err != nil {
		log.Fatal().Err(err).Send()
	}

	printBuildInfo()
	config.Print()
	metricsrv.Serve()

	a, err := app.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = a.Run(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msgf("Shutting down gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err = a.Shutdown(shutdownCtx); err != nil {
		log.Err(err).Send()
	}

	metricsrv.Shutdown(shutdownCtx)
}

func printBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Warn().Msg("No build info")

		return
	}

	fields := make(map[string]any)
	for _, setting := range info.Settings {
		if strings.Contains(setting.Key, "vcs") {
			fields[setting.Key] = setting.Value
		}
	}

	log.Info().Fields(fields).Msg("Build info")
}
