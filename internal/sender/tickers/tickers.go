package tickers

import (
	"time"
)

func StartPoll(interval time.Duration, buffer TickerBuffer) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		buffer.Write()
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
