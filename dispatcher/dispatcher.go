// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package dispatcher

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/host"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	host host.Host
}

var Logger *logrus.Entry

func NewDispatcher(host host.Host) *Dispatcher {
	return &Dispatcher{
		host: host,
	}
}

var lastPostedAt time.Time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func (d *Dispatcher) Dispatch(ctx context.Context, logEventsQueue *queue.Queue, force bool) {
	now := time.Now()
	if !logEventsQueue.Empty() && (force || lastPostedAt.Add(1*time.Minute).Before(now)) {
		Logger.Info("[Dispatch] Dispatching", logEventsQueue.Len(), "log events")
		logEntries, _ := logEventsQueue.Get(logEventsQueue.Len())
		metrics := gatherMetrics(logEntries)
		metrics = aggregateMetrics(metrics, now)
		if len(metrics) > 0 {
			metrics = getOSStat(metrics, now)
			if err := d.host.PostMetrics(metrics); err != nil {
				for _, logEntry := range logEntries {
					logEventsQueue.Put(logEntry)
				}
			}
			lastPostedAt = now
		}
	}
}

type platformInitReportRecordMetrics struct {
	DurationMs float64 `json:"durationMs"`
}
type platformInitReportRecord struct {
	Metrics platformInitReportRecordMetrics `json:"metrics"`
}
type platformInitReport struct {
	Record platformInitReportRecord `json:"record"`
	Time   time.Time                `json:"time"`
}

type platformReportRecordMetrics struct {
	BilledDurationMs float64 `json:"billedDurationMs"`
	DurationMs       float64 `json:"durationMs"`
	InitDurationMs   float64 `json:"initDurationMs"`
	MaxMemoryUsedMB  float64 `json:"maxMemoryUsedMB"`
	MemorySizeMB     float64 `json:"memorySizeMB"`
}
type platformReportRecord struct {
	Metrics platformReportRecordMetrics `json:"metrics"`
}
type platformReport struct {
	Record platformReportRecord `json:"record"`
	Time   time.Time            `json:"time"`
}

type platformRuntimeDoneRecordMetrics struct {
	DurationMs    float64 `json:"durationMs"`
	ProducedBytes float64 `json:"producedBytes"`
}
type platformRuntimeDoneRecord struct {
	Metrics platformRuntimeDoneRecordMetrics `json:"metrics"`
}
type platformRuntimeDone struct {
	Record platformRuntimeDoneRecord `json:"record"`
	Time   time.Time                 `json:"time"`
}

func gatherMetrics(logEntries []interface{}) []*mackerel.MetricValue {
	metrics := make([]*mackerel.MetricValue, 0, 8)
	for _, logEntry := range logEntries {
		switch logEntry.(map[string]interface{})["type"] {
		case "platform.initReport":
			s, _ := json.Marshal(logEntry)
			entry := &platformInitReport{}
			if err := json.Unmarshal(s, &entry); err != nil {
				Logger.Warning("Can't unmarshal platform.initReport:", err)
				continue
			}
			Logger.Info("platform.initReport:", entry)
			metrics = append(
				metrics,
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.initReport.duration.duration",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.DurationMs / 1000.0,
				},
			)

		case "platform.report":
			s, _ := json.Marshal(logEntry)
			entry := &platformReport{}
			if err := json.Unmarshal(s, &entry); err != nil {
				Logger.Warning("Can't unmarshal platform.report:", err)
				continue
			}
			Logger.Info("platform.report:", entry)
			metrics = append(
				metrics,
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.report.billedDuration",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.BilledDurationMs / 1000.0,
				},
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.report.duration",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.DurationMs / 1000.0,
				},
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.report.initDuration",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.InitDurationMs / 1000.0,
				},
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.report.maxMemoryUsed",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.MaxMemoryUsedMB * 1024.0 * 1024.0,
				},
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.report.memorySize",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.MemorySizeMB * 1024.0 * 1024.0,
				},
			)

		case "platform.runtimeDone":
			s, _ := json.Marshal(logEntry)
			entry := &platformRuntimeDone{}
			if err := json.Unmarshal(s, &entry); err != nil {
				Logger.Warning("Can't unmarshal platform.runtimeDone:", err)
				continue
			}
			Logger.Info("platform.runtimeDone:", entry)
			metrics = append(
				metrics,
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.runtimeDone.duration",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.DurationMs / 1000.0,
				},
				&mackerel.MetricValue{
					Name:  "custom.lambda.platform.runtimeDone.producedBytes",
					Time:  entry.Time.Unix(),
					Value: entry.Record.Metrics.ProducedBytes,
				},
			)

		default:
			s, _ := json.Marshal(logEntry)
			Logger.Info("logEntry:", string(s))
		}
	}
	return metrics
}

func aggregateMetrics(metrics []*mackerel.MetricValue, now time.Time) []*mackerel.MetricValue {
	collectedMetricsMap := make(map[string][]*mackerel.MetricValue)
	for _, metric := range metrics {
		if collectedMetricsMap[metric.Name] == nil {
			collectedMetricsMap[metric.Name] = make([]*mackerel.MetricValue, 0, 1)
		}
		collectedMetricsMap[metric.Name] = append(collectedMetricsMap[metric.Name], metric)
	}
	aggregatedMetrics := make([]*mackerel.MetricValue, 0, len(collectedMetricsMap)*3)
	if collectedMetricsMap["custom.lambda.platform.initReport.duration.duration"] != nil {
		aggregatedMetrics = append(aggregatedMetrics, collectedMetricsMap["custom.lambda.platform.initReport.duration.duration"][0])
	}
	for metricName, ms := range collectedMetricsMap {
		if metricName == "custom.lambda.platform.initReport.duration.duration" {
			continue
		}
		var sumValue float64 = 0
		var maxValue float64 = ms[0].Value.(float64)
		var minValue float64 = ms[0].Value.(float64)
		for _, metric := range ms {
			sumValue += metric.Value.(float64)
			maxValue = math.Max(maxValue, metric.Value.(float64))
			minValue = math.Min(maxValue, metric.Value.(float64))
		}
		aggregatedMetrics = append(
			aggregatedMetrics,
			&mackerel.MetricValue{
				Name:  metricName + ".avg",
				Time:  now.Unix(),
				Value: sumValue / float64(len(ms)),
			},
			&mackerel.MetricValue{
				Name:  metricName + ".max",
				Time:  now.Unix(),
				Value: maxValue,
			},
			&mackerel.MetricValue{
				Name:  metricName + ".min",
				Time:  now.Unix(),
				Value: minValue,
			},
		)
	}
	return aggregatedMetrics
}

func getOSStat(metrics []*mackerel.MetricValue, now time.Time) []*mackerel.MetricValue {
	loadavgStat, err := loadavg.Get()
	if err != nil {
		Logger.Warning("Failed to get loadavg:", err)
		return metrics
	}

	metrics = append(
		metrics,
		&mackerel.MetricValue{
			Name:  "custom.lambda.osstat.loadavg.loadavg1",
			Time:  now.Unix(),
			Value: loadavgStat.Loadavg1,
		},
		&mackerel.MetricValue{
			Name:  "custom.lambda.osstat.loadavg.loadavg5",
			Time:  now.Unix(),
			Value: loadavgStat.Loadavg5,
		},
		&mackerel.MetricValue{
			Name:  "custom.lambda.osstat.loadavg.loadavg15",
			Time:  now.Unix(),
			Value: loadavgStat.Loadavg15,
		},
	)

	return metrics
}
