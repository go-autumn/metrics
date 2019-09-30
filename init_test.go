package metrics

import (
	"github.com/go-redis/redis"
	"net/http"
	"testing"
)

var client *redis.Client
var errClient *redis.Client
var redisMetric *AutumnRedisMetric
var redisNoLabelsMetric *AutumnRedisMetric
var metricsCollector *AutumnMetricsCollector
var engine *http.ServeMux

func TestMain(m *testing.M) {
	client = redis.NewClient(&redis.Options{
		Network:    "tcp",
		Addr:       ":6480",
		Password:   "",
		DB:         0,
		MaxRetries: 3,
		PoolSize:   5,
	})

	errClient = redis.NewClient(&redis.Options{
		Network:    "tcp",
		Addr:       ":6480",
		Password:   "",
		DB:         0,
		MaxRetries: 3,
		PoolSize:   5,
	})

	subsystem := "test"
	metricsCollector = NewAutumnCollector(subsystem)

	redisMetric = metricsCollector.NewAutumnRedisMetric(client, "redis_test", "type", "code")
	redisNoLabelsMetric = metricsCollector.NewAutumnRedisMetric(client, "redis_no_labels_test")

	engine = new(http.ServeMux)
	metricsCollector.Use(engine)

	m.Run()
}
