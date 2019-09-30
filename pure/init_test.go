package pure

import (
	"fmt"
	"net/http"
	"testing"
)

var localMetrics *LocalPureMetric
var localNoLabelsMetrics *LocalPureMetric
var metricsCollector *CollectorPure
var engine *http.ServeMux

func TestMain(m *testing.M) {
	subsystem := "test"
	metricsCollector = NewPureCollector(subsystem)
	metricsCollector.UseBeforeHandlerHook(func(ctx *Context) {
		fmt.Println("before handler 1")
	})
	metricsCollector.UseAfterHandlerHook(func(ctx *Context) {
		fmt.Println("after handler 1")
	})

	localMetrics = metricsCollector.NewLocalPureMetric("redis_no_labels_test", "local_type")
	localNoLabelsMetrics = metricsCollector.NewLocalPureMetric("local_no_labels_test")
	localNolLabelValues := metricsCollector.NewLocalPureMetric("test", "local_test1", "type", "code")

	metricsCollector.Register(metricsCollector.NewLocalPureMetric("unregister"))
	metricsCollector.Unregister("unregister")

	metricsCollector.Register(localMetrics, localNoLabelsMetrics, localNolLabelValues)

	metricsCollector.GetMetric("redis_test")
	metricsCollector.GetAllMetrics()

	engine = new(http.ServeMux)
	metricsCollector.Use(engine)
	engine.Handle(DefaultMetricPath, metricsCollector.Handler())
	m.Run()
}
