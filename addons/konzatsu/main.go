package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const pattern = "congestion-thresshold([1-3])"

var (
	url    = flag.String("url", "https://information.konamisportsclub.jp/ksc/004479/", "URL to scrape for congestion level")
	listen = flag.String("listen", ":2112", "Address the metrics server listens on")
	client = &http.Client{}
	regex  *regexp.Regexp
)

func getCongestionLevel() (int, error) {
	resp, err := client.Get(*url)
	if err != nil {
		return 0, fmt.Errorf("could not fetch URL %s: %w", *url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Could not read response data: %w", err)
	}

	matches := regex.FindStringSubmatch(string(data))
	if len(matches) != 2 {
		return 0, fmt.Errorf("regex matched unexpected groups: %v", matches)
	}

	congestionLevel, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("Could not parse congestionLevel: %w", err)
	}

	return congestionLevel, nil
}

type thresholdCollector struct {
	desc *prometheus.Desc
}

func newThresholdCollector() *thresholdCollector {
	return &thresholdCollector{
		desc: prometheus.NewDesc(
			"congestion_threshold",
			"Current congestion threshold level.",
			nil, nil,
		),
	}
}

func (c *thresholdCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *thresholdCollector) Collect(ch chan<- prometheus.Metric) {
	v, err := getCongestionLevel()
	if err != nil {
		// Surfaces the error on the scrape instead of exporting a bogus value
		ch <- prometheus.NewInvalidMetric(c.desc, err)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(v))
}

func main() {
	flag.Parse()

	reg := prometheus.NewRegistry()
	reg.MustRegister(newThresholdCollector())

	regex = regexp.MustCompile(pattern)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	fmt.Printf("Ready, listening on %s\n", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		panic(err)
	}
}
