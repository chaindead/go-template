package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	systemExample = "example"
)

var ExampleCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: systemExample,
		Name:      "count",
		Help:      "The number of inc",
	},
	[]string{"label"},
)
