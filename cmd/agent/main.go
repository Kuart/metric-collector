package main

import (
	agentConfig "github.com/Kuart/metric-collector/config/agent"
	"github.com/Kuart/metric-collector/internal/encryption"
	"github.com/Kuart/metric-collector/internal/sender"
	"github.com/Kuart/metric-collector/internal/sender/buffer"
	"github.com/Kuart/metric-collector/internal/sender/tickers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := agentConfig.New()
	crypto := encryption.New(config.Key)
	client := sender.NewMetricClient(config, crypto)
	bufferToSend := buffer.New()

	osSign := make(chan os.Signal, 1)
	signal.Notify(osSign, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go tickers.StartPoll(config.PollInterval, &bufferToSend)
	go tickers.StartReport(client, config.ReportInterval, &bufferToSend)

	<-osSign
	os.Exit(0)
}
