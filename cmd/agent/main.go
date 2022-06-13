package main

import (
	"github.com/Kuart/metric-collector/cmd"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/sender"
	"github.com/Kuart/metric-collector/internal/storage"
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
	client := sender.NewMetricClient(env.Host, env.Port, pollInterval)
	osSign := make(chan os.Signal, 1)
	signal.Notify(osSign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
			client.SendMetrics(gaugeMetrics, counter)

			counter.Clear()
			gaugeMetrics = *storage.CreateGauge()
		case <-osSign:
			os.Exit(0)
		}
	}
}
