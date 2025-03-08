package app

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	cfg "github.com/spf13/pflag"
)

var (
	setting = cfg.String("domain.setting", "default", "Example for cfg domain with name setting")
)

type App struct {
	setting string
}

func New() (*App, error) {
	return &App{
		setting: *setting,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	log.Info().Msg("App started")

	app.mockWork(ctx)

	return nil
}

func (app *App) Shutdown(_ context.Context) error {
	log.Info().Msg("App stopped")

	return nil
}

func (app *App) mockWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Exit loop on ctx")

			return
		case <-time.After(5 * time.Second):
		}

		log.Info().Str("data", app.setting).Msg("Heartbeat")
	}
}
