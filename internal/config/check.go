package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func checkKeys() {
	yamlKeys := flatten(viper.AllSettings())

	flagKeys := make(map[string]struct{})
	pflag.VisitAll(func(f *pflag.Flag) {
		flagKeys[f.Name] = struct{}{}
	})

	var unknown []string
	for _, key := range yamlKeys {
		if _, ok := flagKeys[key]; !ok {
			unknown = append(unknown, key)
		}
	}

	if len(unknown) > 0 {
		log.Warn().Strs("unknown", unknown).Msg("Found fields in config file, that not defined in code")
	}
}

// flatten returns all keys from a settings map in dot notation.
func flatten(settings map[string]any) []string {
	var keys []string
	flattenDFS(settings, "", &keys)

	return keys
}

func flattenDFS(m map[string]any, prefix string, keys *[]string) {
	for k, v := range m {
		newKey := k
		if prefix != "" {
			newKey = prefix + "." + k
		}

		if subMap, ok := v.(map[string]interface{}); ok {
			flattenDFS(subMap, newKey, keys)

			continue
		}

		*keys = append(*keys, newKey)
	}
}
