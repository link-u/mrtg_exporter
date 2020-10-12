package mrtg

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const namespace = "mrtg_traffic"

var (
	currentIn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "current_in"),
		"Incoming Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	currentOut = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "current_out"),
		"Outgoing Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	averageIn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "average_in"),
		"Average Incoming Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	averageOut = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "average_out"),
		"Average Outgoing Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	maxIn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "max_in"),
		"Max Incoming Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	maxOut = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "max_out"),
		"Max Outgoing Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	averageMaxIn = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "average_max_in"),
		"Average Max Incoming Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
	averageMaxOut = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "average_max_out"),
		"Average Max Outgoing Traffic in Bits per Second",
		[]string{"name", "url", "interval"}, nil,
	)
)

// Describe metrics
func (exp *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- currentIn
	ch <- currentOut
	ch <- maxIn
	ch <- maxOut
	ch <- averageIn
	ch <- averageOut
	ch <- averageMaxIn
	ch <- averageMaxOut
}

// Collect metrics
func (exp *Exporter) Collect(ch chan<- prometheus.Metric) {
	log := zap.L()
	if exp.URL == nil || exp.URL.String() == "" {
		log.Warn("Empty URL", zap.String("url", fmt.Sprintf("%v", exp.URL)))
		return
	}
	name := exp.Name
	var scrapingURL string
	{
		// Create URL without userinfo
		url := *exp.URL
		url.User = nil
		scrapingURL = url.String()
	}
	result, err := exp.Scrape(context.Background())
	if err != nil {
		log.Error("Failed scrape", zap.String("url", exp.URL.String()), zap.Error(err))
		return
	}
	pairs := []struct {
		interval string
		metric   *Metric
	}{
		{
			interval: "daily",
			metric:   &result.Daily,
		},
		{
			interval: "weekly",
			metric:   &result.Weekly,
		},
		{
			interval: "monthly",
			metric:   &result.Monthly,
		},
		{
			interval: "yearly",
			metric:   &result.Yearly,
		},
	}
	for _, pair := range pairs {
		ch <- prometheus.MustNewConstMetric(
			currentIn, prometheus.GaugeValue, float64(pair.metric.CurrentIn), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			currentOut, prometheus.GaugeValue, float64(pair.metric.CurrentOut), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			maxIn, prometheus.GaugeValue, float64(pair.metric.MaxIn), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			maxOut, prometheus.GaugeValue, float64(pair.metric.MaxOut), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			averageIn, prometheus.GaugeValue, float64(pair.metric.AverageIn), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			averageOut, prometheus.GaugeValue, float64(pair.metric.AverageOut), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			averageMaxIn, prometheus.GaugeValue, float64(pair.metric.AverageMaxIn), name, scrapingURL, pair.interval,
		)
		ch <- prometheus.MustNewConstMetric(
			averageMaxOut, prometheus.GaugeValue, float64(pair.metric.AverageMaxOut), name, scrapingURL, pair.interval,
		)
	}
}
