package main

import (
	agentConfig "github.com/Kuart/metric-collector/config/agent"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/sender"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := agentConfig.New()
	client := sender.NewMetricClient(config.Address, *config.PollInterval)

	osSign := make(chan os.Signal, 1)
	signal.Notify(osSign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	pollTicker := time.NewTicker(*config.PollInterval)
	reportTicker := time.NewTicker(*config.ReportInterval)

	var randomGauge metric.Gauge
	counter := metric.GetCounter(0)
	gaugeMetrics := make(map[string]metric.GaugeValue)

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
			gaugeMetrics = make(map[string]metric.GaugeValue)
		case <-osSign:
			os.Exit(0)
		}
	}
}
