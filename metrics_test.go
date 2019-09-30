package metrics

import (
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestMetricsHandler(t *testing.T) {
	_ = Get("/test", nil)
	_ = Get("/metrics", nil)
}

func TestCron(t *testing.T) {
	if err := redisMetric.ResetByCron("* * * * * *"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}

func TestRedisMetrics(t *testing.T) {
	redisNoLabelsMetric.GetLabelValues()
	redisNoLabelsMetric.GetLabels()
	redisNoLabelsMetric.GetValue()
	redisNoLabelsMetric.GetName()
	redisNoLabelsMetric.GetMetricValueKey()
	redisNoLabelsMetric.GetRedisKey()
}

func TestResetRedisKey(t *testing.T) {
	redisNoLabelsMetric.ResetMetricRedisKey("redis_key_reset")
	redisNoLabelsMetric.ResetMetricValueRedisKey("value_key_reset")
}

func Get(uri string, param url.Values) []byte {
	if param != nil && len(param) > 0 {
		uri += "?" + param.Encode()
	}
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	req.Header.Set("test", "1")
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	engine.ServeHTTP(w, req)
	// 提取响应
	result := w.Result()
	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	if body := result.Body; body != nil {
		_ = body.Close()
	}
	return body
}
