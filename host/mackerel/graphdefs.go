package mackerel

import (
	"github.com/mackerelio/mackerel-client-go"
)

var GraphDefs []*mackerel.GraphDefsParam = []*mackerel.GraphDefsParam{
	{
		Name:        "custom.lambda.platform.initReport.duration",
		DisplayName: "Init Duration",
		Unit:        "seconds",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.initReport.duration.duration", DisplayName: "duration", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.report.billedDuration",
		DisplayName: "Billed Duration",
		Unit:        "seconds",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.report.billedDuration.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.report.billedDuration.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.report.billedDuration.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.report.duration",
		DisplayName: "Invoke Duration",
		Unit:        "seconds",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.report.duration.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.report.duration.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.report.duration.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.report.initDuration",
		DisplayName: "Invoke Init Duration",
		Unit:        "seconds",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.report.initDuration.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.report.initDuration.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.report.initDuration.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.report.maxMemoryUsed",
		DisplayName: "Max Memory Used",
		Unit:        "bytes",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.report.maxMemoryUsed.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.report.maxMemoryUsed.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.report.maxMemoryUsed.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.report.memorySize",
		DisplayName: "Memory Size",
		Unit:        "bytes",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.report.memorySize.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.report.memorySize.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.report.memorySize.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.runtimeDone.duration",
		DisplayName: "Done Duration",
		Unit:        "seconds",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.runtimeDone.duration.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.runtimeDone.duration.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.runtimeDone.duration.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.platform.runtimeDone.producedBytes",
		DisplayName: "Produced Bytes",
		Unit:        "bytes",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.platform.runtimeDone.producedBytes.avg", DisplayName: "avg", IsStacked: false},
			{Name: "custom.lambda.platform.runtimeDone.producedBytes.max", DisplayName: "max", IsStacked: false},
			{Name: "custom.lambda.platform.runtimeDone.producedBytes.min", DisplayName: "min", IsStacked: false},
		},
	},
	{
		Name:        "custom.lambda.osstat.loadavg",
		DisplayName: "loadavg",
		Unit:        "float",
		Metrics: []*mackerel.GraphDefsMetric{
			{Name: "custom.lambda.osstat.loadavg.loadavg1", DisplayName: "loadavg1", IsStacked: false},
			{Name: "custom.lambda.osstat.loadavg.loadavg5", DisplayName: "loadavg5", IsStacked: false},
			{Name: "custom.lambda.osstat.loadavg.loadavg15", DisplayName: "loadavg15", IsStacked: false},
		},
	},
}
