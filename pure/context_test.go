package pure

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"testing"
	"time"
)

var testCtx *Context
var middleware http.Handler
var collector MetricCollection
var standardRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "server_request_duration_seconds",
		Help:    "How much time costs(s)",
		Buckets: []float64{0.0005, 0.001, 0.002, 0.005, 0.010, 0.020, 0.050, 0.1, 0.5, 1, 5},
	}, []string{"code", "router"})

func TestInitContext(t *testing.T) {
	testCtx = &Context{
		preHandlers:   make([]HandlerFunc, 0),
		afterHandlers: make([]HandlerFunc, 0),
	}

	prometheus.Register(standardRequestDuration)
	collector = NewPureCollector("test")
	collector.UseBeforeHandlerHook(func(ctx *Context) {
		fmt.Println("this is before 1 start")
		start := time.Now()
		ctx.Set("start", start)
		fmt.Println("this is before 1 end")
	})

	collector.UseAfterHandlerHook(func(ctx *Context) {
		fmt.Println("this is after 1 start")
		start := ctx.GetTime("start")
		elapsed := time.Since(start).Seconds()
		standardRequestDuration.WithLabelValues(fmt.Sprint(ctx.ResponseWriter.Status()), ctx.Request.URL.Path).Observe(elapsed)
		fmt.Println("this is after 1 end")
	})
	middleware = collector.Use(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
	}))
}

func TestContext_Abort(t *testing.T) {
	beforeAbortCtx := &Context{}
	beforeAbortCtx.Abort()

	afterAbortCtx := &Context{handlerType: afterHandler}
	afterAbortCtx.Abort()
}

func TestContext_Next(t *testing.T) {
	beforeCtx := &Context{
		preHandlers:   make([]HandlerFunc, 0),
		afterHandlers: make([]HandlerFunc, 0),
	}
	beforeCtx.preHandlers = append(beforeCtx.preHandlers, func(ctx *Context) {
		fmt.Println("this is before 2 start")
		fmt.Println("this is before 2 end")
	})
	beforeCtx.preHandlers = append(beforeCtx.preHandlers, func(ctx *Context) {
		fmt.Println("this is before 3 start")
		ctx.Next()
		fmt.Println("this is before 3 end")
	})
	beforeCtx.Next()

	afterCtx := &Context{
		handlerType:   afterHandler,
		preHandlers:   make([]HandlerFunc, 0),
		afterHandlers: make([]HandlerFunc, 0),
	}
	afterCtx.afterHandlers = append(afterCtx.afterHandlers, func(ctx *Context) {
		fmt.Println("this is after 2 start")
		fmt.Println("this is after 2 end")
	})
	afterCtx.afterHandlers = append(afterCtx.afterHandlers, func(ctx *Context) {
		fmt.Println("this is after 3 start")
		ctx.Next()
		fmt.Println("this is after 3 end")
	})
	afterCtx.Next()

	errHandlerTypeCtx := &Context{
		handlerType: 2,
	}
	errHandlerTypeCtx.Next()
}

func TestContext_Get(t *testing.T) {
	testCtx.Get("test")
}

func TestContext_MustGet(t *testing.T) {
	testCtx.Set("TestContext_MustGet", "1")
	testCtx.MustGet("TestContext_MustGet")
}

func TestContext_MustGetErr(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
			return
		}
		t.Fatal()
	}()
	testCtx.MustGet("TestContext_MustGetErr")
}

func TestContext_GetString(t *testing.T) {
	testCtx.Set("TestContext_GetString", "1")
	testCtx.GetString("TestContext_GetString")
}

func TestContext_GetBool(t *testing.T) {
	testCtx.Set("TestContext_GetBool", true)
	testCtx.GetBool("TestContext_GetBool")
}

func TestContext_GetInt(t *testing.T) {
	testCtx.Set("TestContext_GetInt", 1)
	testCtx.GetInt("TestContext_GetInt")
}

func TestContext_GetInt64(t *testing.T) {
	testCtx.Set("TestContext_GetInt64", 1)
	testCtx.GetInt64("TestContext_GetInt64")
}

func TestContext_GetFloat64(t *testing.T) {
	testCtx.Set("TestContext_GetFloat64", 1.0)
	testCtx.GetFloat64("TestContext_GetFloat64")
}

func TestContext_GetTime(t *testing.T) {
	testCtx.Set("TestContext_GetTime", time.Now())
	testCtx.GetTime("TestContext_GetTime")
}

func TestContext_GetDuration(t *testing.T) {
	testCtx.Set("TestContext_GetDuration", time.Second)
	testCtx.GetDuration("TestContext_GetDuration")
}

func TestContext_GetStringSlice(t *testing.T) {
	testCtx.Set("TestContext_GetStringSlice", []string{"1"})
	testCtx.GetStringSlice("TestContext_GetStringSlice")
}

func TestContext_GetStringMap(t *testing.T) {
	testCtx.Set("TestContext_GetStringMap", map[string]interface{}{"k": "v"})
	testCtx.GetStringMap("TestContext_GetStringMap")
}

func TestContext_GetStringMapString(t *testing.T) {
	testCtx.Set("TestContext_GetStringMapString", map[string]string{"k": "v"})
	testCtx.GetStringMapString("TestContext_GetStringMapString")
}

func TestContext_GetStringMapStringSlice(t *testing.T) {
	testCtx.Set("TestContext_GetStringMapStringSlice", map[string][]string{"k": {"v"}})
	testCtx.GetStringMapStringSlice("TestContext_GetStringMapStringSlice")
}
