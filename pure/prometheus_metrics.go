package pure

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricType int

const (
	CounterVec MetricType = iota
	Counter
	GaugeVec
	Gauge
	HistogramVec
	Histogram
	SummaryVec
	Summary
)

// NewMetric associates prometheus.Collector based on Metric.Type
func (cp *CollectorPure) NewMetric(name string, metricType MetricType, labels ...string) prometheus.Collector {
	var metric prometheus.Collector
	switch metricType {
	case CounterVec:
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
			labels,
		)
	case Counter:
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
		)
	case GaugeVec:
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
			labels,
		)
	case Gauge:
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
		)
	case HistogramVec:
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
			labels,
		)
	case Histogram:
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
		)
	case SummaryVec:
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
			labels,
		)
	case Summary:
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: cp.subsystem,
				Name:      name,
			},
		)
	}
	return metric
}
