package metrics

import (
	"testing"
)

func TestCollectorIncr(t *testing.T) {
	_ = redisMetric.Incr("1", "300")
	_ = redisNoLabelsMetric.Incr()
	_ = Get("/metrics", nil)
}

func TestCollectorSet(t *testing.T) {
	_ = redisMetric.Set(100, "1", "400")
	_ = redisNoLabelsMetric.Set(101)
	_ = Get("/metrics", nil)
}

func TestCollectorAdd(t *testing.T) {
	_ = redisMetric.Add(1, "1", "200")
	_ = redisNoLabelsMetric.Add(11)
	_ = Get("/metrics", nil)
}
