package pure

import "testing"

func TestNewMetricCounterVec(t *testing.T) {
	metricsCollector.NewMetric("test", CounterVec, "test_counter_vec")
}

func TestNewMetricCounter(t *testing.T) {
	metricsCollector.NewMetric("test", Counter, "test_counter")
}

func TestNewMetricGaugeVec(t *testing.T) {
	metricsCollector.NewMetric("test", GaugeVec, "test_gauge_vec")
}

func TestNewMetricGauge(t *testing.T) {
	metricsCollector.NewMetric("test", Gauge, "test_gauge")
}

func TestNewMetricHistogramVec(t *testing.T) {
	metricsCollector.NewMetric("test", HistogramVec, "test_histogram_vec")
}

func TestNewMetricHistogram(t *testing.T) {
	metricsCollector.NewMetric("test", Histogram, "test_histogram")
}

func TestNewMetricSummaryVec(t *testing.T) {
	metricsCollector.NewMetric("test", SummaryVec, "test_summary_vec")
}

func TestNewMetricSummary(t *testing.T) {
	metricsCollector.NewMetric("test", Summary, "test_summary")
}
