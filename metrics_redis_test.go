package metrics

import (
	"testing"
	"time"
)

func TestNewAutumnRedisMetricMetric(t *testing.T) {
	metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
}

func TestNewAutumnRedis_GetID(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.GetName()
}

func TestNewAutumnRedis_GetLabels(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.GetLabels()
}

func TestNewAutumnRedis_GetLabelValues(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.GetLabelValues()
}

func TestNewAutumnRedis_GetValue(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.Incr("1")
	m.GetValue("1")
}

func TestNewAutumnRedis_Set(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.Set(1)
	m.Set(1, "1")
	m.Set(1, "1")
}

func TestNewAutumnRedis_Incr(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.Incr()
	m.Incr("1")
	m.Incr("1")
}

func TestNewAutumnRedis_Add(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.Add(1)
	m.Add(1, "1")
	m.Add(1, "1")
}

func TestAutumnRedisMetric_Err(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(errClient, "redis_autumn_metrics_test", "type")
	m.Add(1, "1")
	m.Set(1, "1")
	m.Incr("1")

	m2 := metricsCollector.NewAutumnRedisMetric(errClient, "redis_autumn_metrics_test")
	m2.Add(1)
	m2.Set(1)
	m2.Incr()
}

func TestNewAutumnRedis_ResetByCron(t *testing.T) {
	m := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m.Add(1, "1")
	m.ResetByCron("* * * * * *")
	m.ResetByCron("* * * * * *")
	m.ResetByCron("12131")

	m2 := metricsCollector.NewAutumnRedisMetric(client, "redis_autumn_metrics_test", "type")
	m2.Add(1)
	m2.ResetByCron("* * * * * *")
	time.Sleep(time.Second)
}
