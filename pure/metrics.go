// Author SunJun <i@sjis.me>
package pure

// Metrics object interface,
// providing metrics metric data operations and obtaining metrics basic information
type AutumnMetric interface {
	// Get the metrics object id
	GetName() string
	// Get the metrics indicator value
	GetValue(...string) float64
	// Get the metrics parameter set
	GetLabels() []string
	// Get the corresponding values ​​of different parameters of metrics
	GetLabelValues() []*LabelValue
	Set(float64, ...string) error
	Incr(...string) error
	Add(float64, ...string) error
}

// Metrics labels correspond to different combinations of value values
type LabelValue struct {
	Value       float64
	LabelsValue []string
}
