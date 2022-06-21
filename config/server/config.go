package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

var (
	address       = "127.0.0.1:8080"
	restore       = true
	storeInterval = 300 * time.Second
	storeFile     = "/tmp/devops-metrics-db.json"
)

type Config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

type flagConfig struct {
	Address       string
	Restore       bool
	StoreInterval time.Duration
	StoreFile     string
}

func New() Config {
	cfg := Config{}
	flagCfg := flagConfig{}

	flag.StringVar(&flagCfg.Address, "a", address, "server address")
	flag.BoolVar(&flagCfg.Restore, "r", restore, "restore")
	flag.DurationVar(&flagCfg.StoreInterval, "i", storeInterval, "store interval")
	flag.StringVar(&flagCfg.StoreFile, "f", storeFile, "store file path")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if flagCfg.Address != address && cfg.Address == address {
		cfg.Address = flagCfg.Address
	}

	if flagCfg.Restore != restore && cfg.Restore == restore {
		cfg.Restore = flagCfg.Restore
	}

	if flagCfg.StoreInterval != storeInterval && cfg.StoreInterval == storeInterval {
		cfg.StoreInterval = flagCfg.StoreInterval
	}

	if flagCfg.StoreFile != storeFile && cfg.StoreFile == storeFile {
		cfg.StoreFile = flagCfg.StoreFile
	}

	return cfg
}
