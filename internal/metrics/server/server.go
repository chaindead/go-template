package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	cfg "github.com/spf13/pflag"
)

var (
	port = cfg.Int("metrics.port", 9090, "Metric port to listen on")
)

var (
	server *http.Server
)

func Serve() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", *port)
	//nolint: gosec
	server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		url := fmt.Sprintf("http://localhost:%d/metrics", *port)
		log.Info().Str("url", url).Msg("Starting metrics server")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Send()
		}
	}()
}

func Shutdown(ctx context.Context) {
	if server == nil {
		log.Error().Msg("Metrics server is nil")

		return
	}

	err := server.Shutdown(ctx)
	if err != nil {
		log.Err(err).Msgf("Metrics server close")

		return
	}

	log.Info().Msg("Stopped metrics server")
}
