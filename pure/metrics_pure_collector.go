package pure

import (
	"fmt"
	"github.com/go-autumn/lib-metrics/util"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prometheusClient "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK

	contentTypeHeader     = "Content-Type"
	contentEncodingHeader = "Content-Encoding"
	acceptEncodingHeader  = "Accept-Encoding"
)

var (
	defaultMetricsOnce = &sync.Once{}
	reqDurMetric       prometheus.Collector
	reqCntMetric       prometheus.Collector
	reqSzMetric        prometheus.Collector
	resSzMetric        prometheus.Collector

	DefaultMetricPath = "/metrics"
)

var _ MetricCollection = (*CollectorPure)(nil)

type CollectorPure struct {
	locker             *sync.Mutex
	localIp            string
	subsystem          string
	metricsPath        string
	metrics            map[string]AutumnMetric
	prometheusMetrics  []prometheus.Collector
	beforeHandlerHooks []HandlerFunc
	afterHandlerHooks  []HandlerFunc
}

type ResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func NewPureCollector(subsystem string, path ...string) *CollectorPure {
	ips := util.IntranetIP()
	collector := &CollectorPure{
		locker:             new(sync.Mutex),
		localIp:            ips[0],
		subsystem:          subsystem,
		metricsPath:        DefaultMetricPath,
		metrics:            make(map[string]AutumnMetric),
		prometheusMetrics:  make([]prometheus.Collector, 0),
		beforeHandlerHooks: make([]HandlerFunc, 0),
		afterHandlerHooks:  make([]HandlerFunc, 0),
	}
	if len(path) > 0 {
		collector.metricsPath = path[0]
	}
	return collector
}

//  Complete metrics processor with default metrics in the prometheus library
func (cp *CollectorPure) MetricsHandler() http.HandlerFunc {
	register := prometheus.DefaultRegisterer
	gatherer := prometheus.DefaultGatherer
	h := promhttp.InstrumentMetricHandler(
		register, promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{DisableCompression: true}),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		cp.pureMetricsHandler(w, r)
	}
}

//  Complete metrics processor with default metrics in the prometheus library
func (cp *CollectorPure) Handler() http.HandlerFunc {
	return cp.MetricsHandler()
}

// Pure version, only this definition metrics
func (cp *CollectorPure) pureMetricsHandler(w http.ResponseWriter, r *http.Request) {
	gw := io.Writer(w)

	contentType := expfmt.Negotiate(r.Header)
	w.Header().Set(contentTypeHeader, string(contentType))

	// Write custom metrics metrics
	enc := expfmt.NewEncoder(gw, contentType)

	for name, metric := range cp.metrics {

		metricFamily := &prometheusClient.MetricFamily{
			Name:   proto.String(fmt.Sprintf("%s_%s", cp.subsystem, name)),
			Type:   prometheusClient.MetricType_GAUGE.Enum(),
			Metric: make([]*prometheusClient.Metric, 0),
		}

		labels := metric.GetLabels()
		labelValues := metric.GetLabelValues()
		if labels != nil && len(labels) > 0 {
			if nil == labelValues || 0 == len(labelValues) {
				continue
			}
			for idx := range labelValues {
				lvs := labelValues[idx]
				if len(labels) != len(lvs.LabelsValue) {
					continue
				}
				metricDto := &prometheusClient.Metric{
					Gauge: &prometheusClient.Gauge{
						Value: proto.Float64(lvs.Value),
					},
				}
				lps := make([]*prometheusClient.LabelPair, 0, len(labels))
				for idx, label := range lvs.LabelsValue {
					lps = append(lps, &prometheusClient.LabelPair{
						Name:  proto.String(labels[idx]),
						Value: proto.String(label),
					})
				}
				metricDto.Label = lps
				metricFamily.Metric = append(metricFamily.Metric, metricDto)
			}
			if err := enc.Encode(metricFamily); err != nil {
				log.Printf("encode metrics family error:%s\n", err)
			}
			continue
		}

		metricFamily.Metric = append(metricFamily.Metric, &prometheusClient.Metric{
			Gauge: &prometheusClient.Gauge{
				Value: proto.Float64(metric.GetValue()),
			},
		})
		if err := enc.Encode(metricFamily); err != nil {
			log.Printf("encode metrics family error:%s\n", err)
		}
	}
}

// Middleware handler
func (cp *CollectorPure) Use(next http.Handler) http.Handler {
	defaultMetricsOnce.Do(func() {
		reqDurMetric = cp.NewMetric("request_duration_seconds", SummaryVec, "code", "method", "host", "url")
		reqCntMetric = cp.NewMetric("requests_total", CounterVec, "code", "method", "host", "url")
		reqSzMetric = cp.NewMetric("request_size_bytes", Summary)
		resSzMetric = cp.NewMetric("response_size_bytes", Summary)
		prometheus.MustRegister(reqDurMetric, reqCntMetric, reqSzMetric, resSzMetric)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &ResponseWriter{
			ResponseWriter: w,
			size:           noWritten,
			status:         defaultStatus,
		}

		if "/ping" == r.URL.Path || cp.metricsPath == r.URL.Path {
			next.ServeHTTP(sw, r)
			return
		}
		ctx := &Context{
			Request:        r,
			ResponseWriter: sw,
			handlerType:    preHandler,
			preHandlers:    cp.beforeHandlerHooks,
			afterHandlers:  cp.afterHandlerHooks,
		}

		for int(ctx.index) < len(ctx.preHandlers) {
			ctx.preHandlers[ctx.index](ctx)
			ctx.index++
		}

		start := time.Now()
		reqSz := computeApproximateRequestSize(r)

		next.ServeHTTP(sw, r)

		status := strconv.Itoa(sw.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(sw.Size())

		reqDurMetric.(*prometheus.SummaryVec).WithLabelValues(status, r.Method, r.Host, r.URL.Path).Observe(elapsed)
		reqCntMetric.(*prometheus.CounterVec).WithLabelValues(status, r.Method, r.Host, r.URL.Path).Inc()
		reqSzMetric.(prometheus.Summary).Observe(float64(reqSz))
		resSzMetric.(prometheus.Summary).Observe(resSz)

		if !ctx.abort {
			ctx.index = 0
			ctx.handlerType = afterHandler
			for int(ctx.index) < len(ctx.afterHandlers) {
				ctx.afterHandlers[ctx.index](ctx)
				ctx.index++
			}
		}
	})
}

// start with server
func (cp *CollectorPure) ListenAndServe(addr string, h func(http.ResponseWriter, *http.Request)) error {
	http.HandleFunc(cp.metricsPath, h)
	return http.ListenAndServe(addr, nil)
}

func (cp *CollectorPure) Register(ams ...AutumnMetric) {
	cp.locker.Lock()
	for _, am := range ams {
		name := am.GetName()
		if _, ok := cp.metrics[name]; !ok {
			cp.metrics[name] = am
		}
	}
	cp.locker.Unlock()
}

func (cp *CollectorPure) RegisterPrometheus(pmc ...prometheus.Collector) {
	cp.locker.Lock()
	prometheus.MustRegister(pmc...)
	cp.prometheusMetrics = append(cp.prometheusMetrics, pmc...)
	cp.locker.Unlock()
}

func (cp *CollectorPure) Unregister(ids ...string) {
	cp.locker.Lock()
	for _, id := range ids {
		delete(cp.metrics, id)
	}
	cp.locker.Unlock()
}

func (cp *CollectorPure) GetAllMetrics() map[string]AutumnMetric {
	cp.locker.Lock()
	metrics := cp.metrics
	cp.locker.Unlock()
	return metrics
}

func (cp *CollectorPure) GetMetric(id string) (AutumnMetric, bool) {
	cp.locker.Lock()
	am, ok := cp.metrics[id]
	cp.locker.Unlock()
	return am, ok
}

func (cp *CollectorPure) UseBeforeHandlerHook(h ...HandlerFunc) {
	cp.locker.Lock()
	cp.beforeHandlerHooks = append(cp.beforeHandlerHooks, h...)
	cp.locker.Unlock()
}

func (cp *CollectorPure) UseAfterHandlerHook(h ...HandlerFunc) {
	cp.locker.Lock()
	cp.afterHandlerHooks = append(cp.afterHandlerHooks, h...)
	cp.locker.Unlock()
}

func (w *ResponseWriter) Write(data []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if statusCode > 0 && statusCode != w.status {
		w.status = statusCode
	}
}

func (w *ResponseWriter) Status() int {
	return w.status
}

func (w *ResponseWriter) Size() int {
	return w.size
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
