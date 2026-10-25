package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	URL     = "https://information.konamisportsclub.jp/ksc/004479/"
	pattern = "congestion-thresshold([1-3])"
)

var (
	client = &http.Client{}
	regex  *regexp.Regexp
)

func getCongestionLevel() (int, error) {
	resp, err := client.Get(URL)
	defer resp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("Could not fetch URL %s: %w", URL, err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Could not read response data: %w", err)
	}

	matches := regex.FindStringSubmatch(string(data))
	if len(matches) != 2 {
		fmt.Sprintf("RegEx matched not expected: %s", matches)
	}
	congestionLevelStr := matches[1]

	congestionLevel, err := strconv.Atoi(congestionLevelStr)
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
	reg := prometheus.NewRegistry()
	reg.MustRegister(newThresholdCollector())

	regex = regexp.MustCompile(pattern)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	fmt.Println("Ready")
	panic(http.ListenAndServe(":2112", nil))
}
