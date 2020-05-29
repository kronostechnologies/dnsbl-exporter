package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type blCollector struct {
	countMetric *prometheus.Desc
}

func newBlCollector() *blCollector {
	return &blCollector{
		countMetric: prometheus.NewDesc(prometheus.BuildFQName("dnsbl", "ip", "list_count"),
			"Shows on how many lists an IP address was found",
			[]string{"ip", "host", "description"}, nil,
		),
	}
}

func (collector *blCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.countMetric
}

func (collector *blCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range getMetrics(config.Addresses) {
		description := strings.Join(metric.Lists, "/")
		ch <- prometheus.MustNewConstMetric(collector.countMetric, prometheus.GaugeValue, float64(metric.ListCount), metric.IpAddress, metric.Hostname, description)
	}
}
