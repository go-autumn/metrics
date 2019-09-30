![Autumn](autumn.png)

# lib-metrics

[![Build Status](https://travis-ci.org/go-autumn/lib-metrics.svg?branch=master)](https://travis-ci.org/go-autumn/lib-metrics)
[![codecov](https://codecov.io/gh/go-autumn/lib-metrics/branch/master/graph/badge.svg)](https://codecov.io/gh/go-autumn/lib-metrics)

[GoDoc](https://godoc.org/github.com/go-autumn/lib-metrics)

# Dependency
```go
github.com/AlecAivazis/survey/v2 v2.0.1 // indirect
github.com/atotto/clipboard v0.1.2 // indirect
github.com/go-redis/redis v6.15.2+incompatible
github.com/golang/protobuf v1.3.1
github.com/onsi/ginkgo v1.8.0 // indirect
github.com/onsi/gomega v1.5.0 // indirect
github.com/prometheus/client_golang v1.0.0
github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
github.com/prometheus/common v0.6.0
gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5
gopkg.in/yaml.v2 v2.2.2 // indirect
```

# Initialization

init global collector gather
```go
subsystem := "test"
metricsCollector = NewAutumnCollector(subsystem)
```

create local metric and registry
```go
localMetrics = metricsCollector.NewAutumnLocalMetric("local_test", "local_type")
// Reset data at zero every day
localMetrics.ResetByCron("0 0 * * *")
metricsCollector.Register(localMetrics)
```

create metric base on redis and registry
```go
redisMetric = metricsCollector.NewAutumnRedisMetric(client, "redis_test", "type")
// reset redis key
redisMetric.ResetMetricRedisKey("reset_redis_key_test")
// get now redis key
redisMetric.GetRedisKey()
// Reset data at zero every day
redisMetric.ResetByCron("0 0 * * *")
// registry metric
metricsCollector.Register(redisMetric)
```

Register HTTP handler middleware to customize metrics for collecting HTTP requests
```go
// Collect interface responses
var standardRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "server_request_duration_seconds",
		Help:    "How much time costs(s)",
		Buckets: []float64{0.0005, 0.001, 0.002, 0.005, 0.010, 0.020, 0.050, 0.1, 0.5, 1, 5},
	}, []string{"code", "router"})

metricsCollector.RegisterPrometheus(standardRequestDuration)

// Pre-interceptor record start time
metricsCollector.UseBeforeHandlerHook(func(ctx *Context) {
	fmt.Println("this is before 1 start")
	start := time.Now()
	ctx.Set("start", start)
	fmt.Println("this is before 1 end")
})

// the post-interceptor calculates the response time and obtains information such as code, which is set by Observer.
metricsCollector.UseAfterHandlerHook(func(ctx *Context) {
	fmt.Println("this is after 1 start")
	start := ctx.GetTime("start")
	elapsed := time.Since(start).Seconds()
	standardRequestDuration.WithLabelValues(fmt.Sprint(ctx.ResponseWriter.Status()), ctx.Request.URL.Path).Observe(elapsed)
	fmt.Println("this is after 1 end")
})
```

# Usage

```go
// local metrics
localMetrics.Set(1,"1")
localMetrics.Incr("2")
localMetrics.Add(1,"3")

// redis metrics
redisMetric.Set(1,"1")
redisMetric.Incr("2")
redisMetric.Add(1,"3")

// Using native metric
standardRequestDuration.WithLabelValues(200, "/test").Observe(10)
```

# Registry Handler

Register gin middleware and add `metrics` interface to collect HTTP metrics information
```go
engine = gin.New()
// add gin middleware
engine.Use(func(ctx *gin.Context) {
	metricsCollector.Use(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		ctx.Next()
	})).ServeHTTP(ctx.Writer,ctx.Request)
})
// Register gin interface to provide metrics interface
engine.Get(DefaultMetricPath, gin.WrapF(metricsCollector.Handler()))
``` 

Register native HTTP Middleware
```go
engine = new(http.ServeMux)
h:=metricsCollector.Use(engine)
engine.Handle(DefaultMetricPath, metricsCollector.Handler())
http.ListenAndServe(":8080", h)
```

Built-in default HTTP service startup metrics interface
```go
if err := metricsCollector.ListenAndServe(":8080",metricsCollector.Handler); err != nil && err != http.ErrServerClosed {
     panic(err)
}
```