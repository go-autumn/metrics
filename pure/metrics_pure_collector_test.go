package pure

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectorPure_RegisterPrometheus(t *testing.T) {
	s := prometheus.NewSummary(prometheus.SummaryOpts{Name: "test_prometheus"})
	collector.RegisterPrometheus(s)
}

func TestNewPureCollector(t *testing.T) {
	localMetric := metricsCollector.NewLocalPureMetric("test_local", "type")
	localMetric.Incr("1")
	metricsCollector.Register(localMetric)
}

func TestNewAutumnCollector(t *testing.T) {
	NewPureCollector("test", "/metrics")
}

func TestAutumnCollector_handler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/test/handler", nil)
	r.Header.Set("k", "v")
	w := httptest.NewRecorder()
	metricsCollector.Handler()(w, r)
}

func TestAutumnCollector_Use(t *testing.T) {
	middleware := metricsCollector.Use(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.Header.Set("k", "v")
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, r)
}

func TestAutumnCollector_ListenAndServe(t *testing.T) {
	go func() {
		_ = metricsCollector.ListenAndServe(":8888", metricsCollector.Handler())
	}()
}

func TestLocalPureMetric(t *testing.T) {
	_ = localMetrics.Incr("1")
	_ = localNoLabelsMetrics.Incr()

	_ = localMetrics.Set(200, "1")
	_ = localNoLabelsMetrics.Set(201)

	_ = localMetrics.Add(1, "1")
	_ = localNoLabelsMetrics.Add(22)
}

func TestResponseWriter_Write(t *testing.T) {
	w := httptest.NewRecorder()
	write := &ResponseWriter{
		ResponseWriter: w,
	}
	write.Write([]byte(""))
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	write := &ResponseWriter{
		ResponseWriter: w,
	}
	write.WriteHeader(http.StatusOK)
}
