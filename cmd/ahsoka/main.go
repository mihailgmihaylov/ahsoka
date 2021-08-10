package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	rc      = flag.Float64("rc", 0, "Set the value that should be pushed as an Pushgateway event.")
	pushUrl = flag.String("push-url", "https://push.staging-cluster.receipt-labs.com", "Set the value that should be pushed as an Pushgateway event.")

	completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bootstrap_completion_timestamp_seconds",
		Help: "The timestamp of the finish of the bootstrap scritpt, successful or not.",
	})
	successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bootstrap_success_timestamp_seconds",
		Help: "The timestamp of the successful finish of the bootstrap scritpt.",
	})
	returnCode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bootstrap_return_code",
		Help: "The return code of the script",
	})
)

func main() {
	flag.Parse()
	registry := prometheus.NewRegistry()
	registry.MustRegister(completionTime, returnCode)

	pusher := push.New(*pushUrl, "bootstrap").Gatherer(registry)

	returnCode.Set(*rc)
	completionTime.SetToCurrentTime()
	if *rc != 0 {
		fmt.Println("Script failed with return code", rc)
	} else {
		pusher.Collector(successTime)
		successTime.SetToCurrentTime()
	}
	if err := pusher.Add(); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
}
