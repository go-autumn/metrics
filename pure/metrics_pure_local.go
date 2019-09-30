package pure

import (
	"errors"
	"github.com/go-autumn/lib-metrics/util"
	"strings"
	"sync"
)

type LocalPureMetric struct {
	Name        string
	Type        MetricType
	Labels      []string
	LabelValues map[string]*LabelValue
	locker      *sync.Mutex
	value       float64
	subsystem   string
	localIp     string
}

var _ AutumnMetric = (*LocalPureMetric)(nil)

func (cp *CollectorPure) NewLocalPureMetric(name string, labels ...string) *LocalPureMetric {
	metric := &LocalPureMetric{
		locker:      new(sync.Mutex),
		subsystem:   cp.subsystem,
		Name:        name,
		Type:        GaugeVec,
		Labels:      labels,
		LabelValues: make(map[string]*LabelValue),
	}
	ips := util.IntranetIP()
	metric.localIp = ips[0]
	return metric
}

func (arm *LocalPureMetric) GetName() string {
	return arm.Name
}

func (arm *LocalPureMetric) GetValue(labelValues ...string) float64 {
	arm.locker.Lock()
	defer arm.locker.Unlock()
	if len(labelValues) > 0 && len(arm.LabelValues) > 0 {
		for k, v := range arm.LabelValues {
			if k == strings.Join(labelValues, ",") {
				return v.Value
			}
		}
	}
	return arm.value
}

func (arm *LocalPureMetric) GetLabels() []string {
	return arm.Labels
}

func (arm *LocalPureMetric) GetLabelValues() []*LabelValue {
	arm.locker.Lock()
	defer arm.locker.Unlock()
	r := make([]*LabelValue, 0, len(arm.LabelValues))
	for _, v := range arm.LabelValues {
		r = append(r, v)
	}
	return r
}

func (arm *LocalPureMetric) Set(value float64, labelValues ...string) error {
	if len(arm.Labels) != len(labelValues) {
		return errors.New("labels key no match label values")
	}
	arm.locker.Lock()
	defer arm.locker.Unlock()

	if len(arm.Labels) > 0 {
		lv := strings.Join(labelValues, ",")
		if _, ok := arm.LabelValues[lv]; ok {
			arm.LabelValues[lv].Value = value
		} else {
			arm.LabelValues[lv] = &LabelValue{
				Value:       value,
				LabelsValue: labelValues,
			}
		}
		return nil
	}
	arm.value = value
	return nil
}

func (arm *LocalPureMetric) Incr(labelValues ...string) error {
	if len(arm.Labels) != len(labelValues) {
		return errors.New("labels key no match label values")
	}
	arm.locker.Lock()
	defer arm.locker.Unlock()

	if len(arm.Labels) > 0 {
		lv := strings.Join(labelValues, ",")
		if _, ok := arm.LabelValues[lv]; ok {
			arm.LabelValues[lv].Value += 1
		} else {
			arm.LabelValues[lv] = &LabelValue{
				Value:       1,
				LabelsValue: labelValues,
			}
		}
		return nil
	}
	arm.value += 1
	return nil
}

func (arm *LocalPureMetric) Add(value float64, labelValues ...string) error {
	if len(arm.Labels) != len(labelValues) {
		return errors.New("labels key no match label values")
	}
	arm.locker.Lock()
	defer arm.locker.Unlock()

	if len(arm.Labels) > 0 {
		lv := strings.Join(labelValues, ",")
		if _, ok := arm.LabelValues[lv]; ok {
			arm.LabelValues[lv].Value += value
		} else {
			arm.LabelValues[lv] = &LabelValue{
				Value:       value,
				LabelsValue: labelValues,
			}
		}
		return nil
	}
	arm.value += value
	return nil
}
