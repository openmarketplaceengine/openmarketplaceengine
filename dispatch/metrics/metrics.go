package metrics

import "github.com/prometheus/client_golang/prometheus"

var MatrixApiHits = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace:   "google",
	Subsystem:   "",
	Name:        "matrix_api_hits",
	Help:        "",
	ConstLabels: nil,
})

var MatrixApiCallDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace:   "google",
	Subsystem:   "",
	Name:        "matrix_api_call_duration",
	Help:        "",
	ConstLabels: nil,
	Buckets:     nil,
})

var EstimatesApiCallDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace:   "demand",
	Subsystem:   "",
	Name:        "estimates_api_call_duration",
	Help:        "",
	ConstLabels: nil,
	Buckets:     []float64{50, 100, 150, 200, 300, 500, 700},
})
