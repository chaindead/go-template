package config

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	path = pflag.StringP("config", "c", "", "Path to config file")
)

func Parse() error {
	pflag.Parse()

	if *path != "" {
		dir, filename, ext := parsePath(*path)

		viper.SetConfigName(filename)
		viper.SetConfigType(ext)
		viper.AddConfigPath(dir)

		err := viper.ReadInConfig()
		if err != nil {
			return errors.Wrapf(err, "read config file(%s)", *path)
		}

		checkKeys()
	}

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return errors.Wrap(err, "bind flags")
	}

	pflag.VisitAll(updateWithViperValue)

	return nil
}

func updateWithViperValue(f *pflag.Flag) {
	value := viper.GetString(f.Name)

	if err := f.Value.Set(value); err != nil {
		log.Warn().Err(err).
			Str("key", f.Name).
			Str("value", value).
			Msgf("Fail to properly set config value")
	}
}

func parsePath(path string) (dir, filename, ext string) {
	dir = filepath.Dir(path)
	base := filepath.Base(path)

	ext = filepath.Ext(base)
	filename = strings.TrimSuffix(base, ext)

	ext = strings.TrimLeft(ext, ".")

	return dir, filename, ext
}

func Print() {
	settings := viper.AllSettings()

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().RawJSON("data", data).Msg("Config")
}
