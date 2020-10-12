package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/link-u/mrtg_exporter/internal/mrtg"
	"github.com/mattn/go-isatty"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write(
		[]byte(fmt.Sprintf(`
<html>
	<head><title>MRTG Exporter</title></head>
	<body>
		<h1>MRTG Exporter</h1>
		<p><a href="%s">Metrics</a></p>
	</body>
</html>
`, *metricsPath)))
}

var probePath = flag.String("web.probe-path", "/probe", "Path under which to expose metrics")

var whiteSpace = regexp.MustCompile(`\s+`)

func probeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	targetRaw := r.URL.Query().Get("target")
	targets := whiteSpace.Split(targetRaw, -1)
	var name string
	var urlRaw string
	var scrapingURL *url.URL
	if len(targets) >= 1 {
		name = targets[0]
		urlRaw = targets[1]
	} else {
		name = targets[0]
		urlRaw = targets[0]
	}
	scrapingURL, err = url.Parse(urlRaw)
	if err != nil {
		http.Error(w, fmt.Sprintf("Inalid URL: %s\nErr=%v", targetRaw, err), 500)
		return
	}
	// The following timeout block was taken wholly from the blackbox exporter
	//   https://github.com/prometheus/blackbox_exporter/blob/master/main.go
	var timeoutSeconds float64 = 0
	if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
		var err error
		timeoutSeconds, err = strconv.ParseFloat(v, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse timeout from Prometheus header: %s", err), http.StatusInternalServerError)
			return
		}
	} else {
		timeoutSeconds = 10
	}
	if timeoutSeconds == 0 {
		timeoutSeconds = 10
	}

	exp := &mrtg.Exporter{}
	exp.Name = name
	exp.Timeout = time.Duration(float64(time.Second) * timeoutSeconds)
	exp.URL = scrapingURL

	registry := prometheus.NewRegistry()
	registry.MustRegister(exp)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

var metricsPath = flag.String("web.metric-path", "/metric", "Path under which to expose metrics")

var listenAddress = flag.String("web.listen-address", ":9230", "Address to listen on for web interface and telemetry.")

func main() {
	var err error
	var log *zap.Logger
	flag.Parse()

	// Check is terminal
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create logger: %v", err)
		os.Exit(-1)
	}
	undo := zap.ReplaceGlobals(log)
	defer undo()
	log.Info("Log System Initialized.")

	http.HandleFunc("/", indexHandler)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc(*probePath, probeHandler)

	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal("Faled to run web server", zap.Error(err))
	}
}
