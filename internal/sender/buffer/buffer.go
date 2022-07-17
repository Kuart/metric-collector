package buffer

import (
	"github.com/Kuart/metric-collector/internal/metric"
	"sync"
)

type BufferToSend struct {
	mu           sync.RWMutex
	counter      metric.Counter
	gaugeMetrics metric.GaugeState
}

func New() BufferToSend {
	return BufferToSend{
		counter:      metric.GetCounter(0),
		gaugeMetrics: make(metric.GaugeState),
	}
}

func (bts *BufferToSend) Write() {
	bts.mu.Lock()
	defer bts.mu.Unlock()

	bts.counter.PollTick()
	randomGauge := metric.GetRandomGauge()
	bts.gaugeMetrics[randomGauge.Name] += randomGauge.Value

	for _, item := range metric.GetGauge() {
		bts.gaugeMetrics[item.Name] += item.Value
	}
}

func (bts *BufferToSend) Clear() {
	bts.mu.Lock()
	defer bts.mu.Unlock()

	bts.counter = metric.GetCounter(0)
	bts.gaugeMetrics = make(metric.GaugeState)
}

func (bts BufferToSend) Get() (metric.GaugeState, metric.Counter) {
	return bts.gaugeMetrics, bts.counter
}
