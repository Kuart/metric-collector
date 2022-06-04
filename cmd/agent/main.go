package main

import (
	"github.com/Kuart/metric-collector/internal/api"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	osSign := make(chan os.Signal, 1)
	signal.Notify(osSign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	client := http.Client{Timeout: pollInterval}
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	var randomGauge metric.Gauge
	counter := metric.GetCounter(0)
	gaugeMetrics := *storage.CreateGauge()

	for {
		select {
		case <-pollTicker.C:
			counter.PollTick()
			randomGauge = metric.GetRandomGauge()
			gaugeMetrics[randomGauge.Name] += randomGauge.Value

			for _, item := range metric.GetGauge() {
				gaugeMetrics[item.Name] += item.Value
			}
		case <-reportTicker.C:
			api.SendMetric(&client, gaugeMetrics, counter)

			counter.Clear()
			gaugeMetrics = *storage.CreateGauge()
		case <-osSign:
			os.Exit(0)
		}
	}
}
