package tickers

import (
	"time"
)

func StartCommon(interval time.Duration, buffer TickerBuffer) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		buffer.WriteCommon()
	}
}

func StartGopsutil(interval time.Duration, buffer TickerBuffer) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		buffer.WriteGopsutil()
	}
}

func StartReport(client HTTPClient, interval time.Duration, buffer TickerBuffer) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		gauge, counter := buffer.Get()
		client.SendMetrics(gauge, counter)
		buffer.Clear()
	}
}
