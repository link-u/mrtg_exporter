package mrtg

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Exporter struct {
	Timeout time.Duration
	URL     *url.URL
}

var metricPattern = regexp.MustCompile(`<!-- (maxin|maxout|avin|avout|cuin|cuout|avmxin|avmxout) ([dwmy]) (\d+) -->`)

func (exp *Exporter) Scrape(c context.Context) (*Metrics, error) {
	var metrics Metrics
	client := &http.Client{Timeout: exp.Timeout}

	req, err := http.NewRequestWithContext(c, "GET", exp.URL.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	for _, v := range metricPattern.FindAllSubmatch(html, -1) {
		var err error
		dataType := string(v[1])
		duration := string(v[2])
		dat := string(v[3])
		var value *uint64
		var metric *Metric
		switch duration {
		case "d":
			metric = &metrics.Daily
		case "w":
			metric = &metrics.Weekly
		case "m":
			metric = &metrics.Monthly
		case "y":
			metric = &metrics.Yearly
		default:
			panic(fmt.Sprintf("Invalid duration: %s", duration))
		}
		switch dataType {
		case "cuin":
			value = &metric.CurrentIn
		case "cuout":
			value = &metric.CurrentOut
		case "maxin":
			value = &metric.MaxIn
		case "maxout":
			value = &metric.MaxOut
		case "avin":
			value = &metric.AverageIn
		case "avout":
			value = &metric.AverageOut
		case "avmxin":
			value = &metric.AverageMaxIn
		case "avmxout":
			value = &metric.AverageMaxOut
		default:
			panic(fmt.Sprintf("Invalid data type: %s", dataType))
		}
		*value, err = strconv.ParseUint(dat, 10, 64)
		if err != nil {
			return nil, err
		}
		*value = *value * 8 // Convert 'Bytes/Sec' to 'Bits/Sec'
	}
	return &metrics, nil
}
