package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

var (
	address        = "127.0.0.1:8080"
	reportInterval = 10 * time.Second
	pollInterval   = 2 * time.Second
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

type flagConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func New() Config {
	cfg := Config{}
	flagCfg := flagConfig{}

	flag.StringVar(&flagCfg.Address, "a", address, "server address")
	flag.DurationVar(&flagCfg.ReportInterval, "r", reportInterval, "report interval")
	flag.DurationVar(&flagCfg.PollInterval, "p", pollInterval, "poll interval")
	flag.Parse()

	log.Printf("Agent init flags: %+v\n", flagCfg)

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	log.Printf("Agent init ENV: %+v\n", cfg)

	if flagCfg.Address != address && cfg.Address == address {
		cfg.Address = flagCfg.Address
	}

	if flagCfg.ReportInterval != reportInterval && cfg.ReportInterval == reportInterval {
		cfg.ReportInterval = flagCfg.ReportInterval
	}

	if flagCfg.PollInterval != pollInterval && cfg.PollInterval == pollInterval {
		cfg.PollInterval = flagCfg.PollInterval
	}

	log.Printf("Agent config: %+v\n", cfg)

	return cfg
}
