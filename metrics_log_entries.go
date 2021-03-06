package main

import (
	"context"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollectorLogEntries struct {
	CollectorProcessorGeneral

	prometheus struct {
		service *prometheus.GaugeVec
	}

	isoverview bool
	timezone   string
	since      string
	until      string
}

func (m *MetricsCollectorLogEntries) Setup(collector *CollectorGeneral) {
	m.CollectorReference = collector

	m.prometheus.service = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pagerduty_log_entries",
			Help: "PagerDuty log entries",
		},
		[]string{
			"entryId",
			"summary",
			"time",
			"activity",
		},
	)

	prometheus.MustRegister(m.prometheus.service)
}

func (m *MetricsCollectorLogEntries) Reset() {
	m.prometheus.service.Reset()
}

func (m *MetricsCollectorLogEntries) Collect(ctx context.Context, callback chan<- func()) {
	listOpts := pagerduty.ListLogEntriesOptions{}
	listOpts.Limit = PagerdutyListLimit
	listOpts.Offset = 0

	if m.since != "" {
		listOpts.Since = m.since
		daemonLogger.Verbosef("Since:  %v", listOpts.Since)
	}
	if m.until != "" {
		listOpts.Until = m.until
		daemonLogger.Verbosef("Until:  %v", listOpts.Until)
	}
	if m.timezone != "" {
		listOpts.TimeZone = m.timezone
	} else {
		listOpts.TimeZone = "UTC"
	}
	daemonLogger.Verbosef("TimeZone:  %v", listOpts.TimeZone)
	if m.isoverview {
		listOpts.IsOverview = true
	}
	daemonLogger.Verbosef("IsOverview:  %v", listOpts.IsOverview)

	logEntriesMetricList := MetricCollectorList{}

	for {
		daemonLogger.Verbosef(" - fetch log entries (offset: %v, limit:%v)", listOpts.Offset, listOpts.Limit)

		list, err := PagerDutyClient.ListLogEntries(listOpts)
		m.CollectorReference.PrometheusAPICounter().WithLabelValues("ListLogEntries").Inc()

		if err != nil {
			panic(err)
		}

		for _, logEntry := range list.LogEntries {
			logEntriesMetricList.AddInfo(prometheus.Labels{
				"entryId":  logEntry.ID,
				"summary":  logEntry.Summary,
				"time":     logEntry.CreatedAt,
				"activity": logEntry.Incident.Summary,
			})
		}

		listOpts.Offset += list.Limit
		if !list.More {
			break
		}
	}

	// set metrics
	callback <- func() {
		logEntriesMetricList.GaugeSet(m.prometheus.service)
	}
}
