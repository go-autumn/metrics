package metrics

import (
	"github.com/go-autumn/lib-metrics/pure"
)

var _ pure.MetricCollection = (*AutumnMetricsCollector)(nil)

type AutumnMetricsCollector struct {
	*pure.CollectorPure
	subsystem   string
	metricsPath string
}

func NewAutumnCollector(subsystem string, path ...string) *AutumnMetricsCollector {
	pureCollector := pure.NewPureCollector(subsystem, path...)
	collector := &AutumnMetricsCollector{
		CollectorPure: pureCollector,
		subsystem:     subsystem,
		metricsPath:   pure.DefaultMetricPath,
	}
	return collector
}
