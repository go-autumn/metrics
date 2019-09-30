package metrics

import (
	"errors"
	"fmt"
	"github.com/go-autumn/lib-metrics/pure"
	"github.com/go-autumn/lib-metrics/util"
	"github.com/go-redis/redis"
	"gopkg.in/robfig/cron.v2"
	"strconv"
	"strings"
	"sync"
)

const (
	DefaultMetricRedisKey  = "a:p:rm:%s:%s:%s"
	DefaultMetricsValueKey = "m_value"
)

type AutumnRedisMetric struct {
	Name                string
	Type                pure.MetricType
	subsystem           string
	localIp             string
	redisKey            string
	metricValueRedisKey string
	Labels              []string
	locker              *sync.Mutex
	client              *redis.Client
	cronEntryID         cron.EntryID
}

var _ pure.AutumnMetric = (*AutumnRedisMetric)(nil)

func (ac *AutumnMetricsCollector) NewAutumnRedisMetric(client *redis.Client, name string, labels ...string) *AutumnRedisMetric {
	metric := &AutumnRedisMetric{
		client:              client,
		locker:              new(sync.Mutex),
		Name:                name,
		Type:                pure.GaugeVec,
		Labels:              labels,
		subsystem:           ac.subsystem,
		metricValueRedisKey: DefaultMetricsValueKey,
	}
	ips := util.IntranetIP()
	metric.localIp = ips[0]
	metric.redisKey = fmt.Sprintf(DefaultMetricRedisKey, ac.subsystem, ips[0], name)
	return metric
}

func (arm *AutumnRedisMetric) ResetMetricRedisKey(redisKey string) *AutumnRedisMetric {
	arm.redisKey = redisKey
	return arm
}

func (arm *AutumnRedisMetric) ResetMetricValueRedisKey(valueKey string) *AutumnRedisMetric {
	arm.metricValueRedisKey = valueKey
	return arm
}

func (arm *AutumnRedisMetric) GetRedisKey() string {
	return arm.redisKey
}

func (arm *AutumnRedisMetric) GetMetricValueKey() string {
	return arm.metricValueRedisKey
}

func (arm *AutumnRedisMetric) GetName() string {
	return arm.Name
}

func (arm *AutumnRedisMetric) GetValue(label ...string) float64 {
	key := arm.metricValueRedisKey
	if len(label) > 0 {
		key = label[0]
	}
	f, _ := arm.client.HGet(arm.redisKey, key).Float64()
	return f
}

func (arm *AutumnRedisMetric) GetLabels() []string {
	return arm.Labels
}

func (arm *AutumnRedisMetric) GetLabelValues() []*pure.LabelValue {
	allLabels := arm.client.HGetAll(arm.redisKey).Val()
	delete(allLabels, arm.metricValueRedisKey)
	r := make([]*pure.LabelValue, 0, len(allLabels))
	for k, vStr := range allLabels {
		v, _ := strconv.ParseFloat(vStr, 64)
		lvs := strings.Split(k, ",")
		r = append(r, &pure.LabelValue{
			Value:       v,
			LabelsValue: lvs,
		})
	}
	return r
}

func (arm *AutumnRedisMetric) Set(value float64, labelValues ...string) error {
	if len(arm.Labels) != len(labelValues) {
		return errors.New("labels key no match label values")
	}
	if len(arm.Labels) > 0 {
		if err := arm.client.HSet(arm.redisKey, strings.Join(labelValues, ","), value).Err(); err != nil {
			return err
		}
		return nil
	}
	if err := arm.client.HSet(arm.redisKey, arm.metricValueRedisKey, value).Err(); err != nil {
		return err
	}
	return nil
}

func (arm *AutumnRedisMetric) Incr(labelValues ...string) error {
	if len(arm.Labels) != len(labelValues) {
		return errors.New("labels key no match label values")
	}
	if len(arm.Labels) > 0 {
		if err := arm.client.HIncrByFloat(arm.redisKey, strings.Join(labelValues, ","), 1).Err(); err != nil {
			return err
		}
		return nil
	}
	if err := arm.client.HIncrByFloat(arm.redisKey, arm.metricValueRedisKey, 1).Err(); err != nil {
		return err
	}
	return nil
}

func (arm *AutumnRedisMetric) Add(value float64, labelValues ...string) error {
	if len(arm.Labels) != len(labelValues) {
		return errors.New("labels key no match label values")
	}
	if len(arm.Labels) > 0 {
		if err := arm.client.HIncrByFloat(arm.redisKey, strings.Join(labelValues, ","), value).Err(); err != nil {
			return err
		}
		return nil
	}
	if err := arm.client.HIncrByFloat(arm.redisKey, arm.metricValueRedisKey, value).Err(); err != nil {
		return err
	}
	return nil
}

func (arm *AutumnRedisMetric) ResetByCron(spec string) error {
	arm.locker.Lock()
	defer arm.locker.Unlock()

	if arm.cronEntryID > 0 {
		DefaultCron().Remove(arm.cronEntryID)
	}
	id, err := DefaultCron().AddFunc(spec, func() {
		labels := arm.GetLabels()
		if labels != nil && len(labels) > 0 {
			labelValues := arm.GetLabelValues()
			for _, lvs := range labelValues {
				arm.Set(0, strings.Join(lvs.LabelsValue, ","))
			}
			return
		}
		arm.Set(0)
	})
	if err != nil {
		return err
	}
	arm.cronEntryID = id
	return nil
}
