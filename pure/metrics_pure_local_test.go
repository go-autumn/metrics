package pure

import (
	"testing"
)

func TestNewAutumnLocalMetric(t *testing.T) {
	metricsCollector.NewLocalPureMetric("test1", "type")
}

func TestAutumnLocalMetric_GetID(t *testing.T) {
	localMetrics.GetName()
}

func TestAutumnLocalMetric_GetLabels(t *testing.T) {
	localMetrics.GetLabels()
}

func TestAutumnLocalMetric_GetLabelValues(t *testing.T) {
	localMetrics.GetLabelValues()
}

func TestAutumnLocalMetric_GetValue(t *testing.T) {
	localMetrics.Incr("1")
	localMetrics.GetValue("1")
}

func TestAutumnLocalMetric_Set(t *testing.T) {
	localMetrics.Set(1)
	localMetrics.Set(1, "1")
	localMetrics.Set(1, "1")
}

func TestAutumnLocalMetric_Incr(t *testing.T) {
	localMetrics.Incr()
	localMetrics.Incr("1")
	localMetrics.Incr("1")
}

func TestAutumnLocalMetric_Add(t *testing.T) {
	localMetrics.Add(1)
	localMetrics.Add(1, "1")
	localMetrics.Add(1, "1")
}
