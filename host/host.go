package host

import "github.com/mackerelio/mackerel-client-go"

type Host interface {
	Retire() error
	CreateGraphDefs() error
	PostMetrics(metrics []*mackerel.MetricValue) error
}
