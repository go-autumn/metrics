package metrics

import (
	"github.com/go-autumn/lib-metrics/pure"
	"gopkg.in/robfig/cron.v2"
	"strings"
	"sync"
)

type AutumnLocalMetric struct {
	*pure.LocalPureMetric
	locker      *sync.Mutex
	cronEntryID cron.EntryID
}

var _ pure.AutumnMetric = (*AutumnLocalMetric)(nil)

func (ac *AutumnMetricsCollector) NewAutumnLocalMetric(name string, labels ...string) *AutumnLocalMetric {
	pureLocalMetrics := ac.NewLocalPureMetric(name, labels...)
	metric := &AutumnLocalMetric{
		LocalPureMetric: pureLocalMetrics,
		locker:          new(sync.Mutex),
	}
	return metric
}

func (arm *AutumnLocalMetric) ResetByCron(spec string) error {
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
