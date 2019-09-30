// Author SunJun <i@sjis.me>
package pure

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

// Metrics collection interface for registering metrics objects,
// providing metrics interfaces and services
type MetricCollection interface {
	// Register metrics
	Register(...AutumnMetric)
	// Register Prometheus native metrics collector
	RegisterPrometheus(...prometheus.Collector)
	// Unregister metrics
	Unregister(...string)
	// Http middleware, Used to collect http indicator information
	Use(http.Handler) http.Handler
	// Get a metrics object
	GetMetric(string) (AutumnMetric, bool)
	// Full metrics interface, return full metrics metrics
	MetricsHandler() http.HandlerFunc
	// Full metrics interface, return full metrics metrics
	Handler() http.HandlerFunc
	// Start the http service and provide the metrics interface
	ListenAndServe(string, func(http.ResponseWriter, *http.Request)) error
	// Add http middleware handler pre hook
	UseBeforeHandlerHook(...HandlerFunc)
	// Add http middleware handler after hook
	UseAfterHandlerHook(...HandlerFunc)
}
