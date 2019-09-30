package metrics

import (
	"testing"
	"time"
)

func TestAutumnLocalMetric_ResetByCron_Error(t *testing.T) {
	m := metricsCollector.NewAutumnLocalMetric("test1")
	m.ResetByCron("123")
}

func TestAutumnLocalMetric_Labels_ResetByCron(t *testing.T) {
	m := metricsCollector.NewAutumnLocalMetric("test2", "type")
	m.Add(1, "1")
	m.ResetByCron("* * * * * *")
	m.ResetByCron("* * * * * *")
	time.Sleep(time.Second)
}

func TestAutumnLocalMetric_No_Labels_ResetByCron(t *testing.T) {
	m := metricsCollector.NewAutumnLocalMetric("test3")
	m.Add(1)
	m.ResetByCron("* * * * * *")
	time.Sleep(time.Second)
}
