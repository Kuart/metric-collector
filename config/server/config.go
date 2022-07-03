package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

var (
	address       = "127.0.0.1:8080"
	restore       = true
	storeInterval = 300 * time.Second
	storeFile     = "/tmp/devops-metrics-db.json"
	key           = ""
)

type Config struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
	Key           string        `env:"KEY" envDefault:"test"`
}

type flagConfig struct {
	Address       string
	Restore       bool
	StoreInterval time.Duration
	StoreFile     string
	Key           string
}

func New() Config {
	cfg := Config{}
	flagCfg := flagConfig{}

	flag.StringVar(&flagCfg.Address, "a", address, "server address")
	flag.BoolVar(&flagCfg.Restore, "r", restore, "restore")
	flag.DurationVar(&flagCfg.StoreInterval, "i", storeInterval, "store interval")
	flag.StringVar(&flagCfg.StoreFile, "f", storeFile, "store file path")
	flag.StringVar(&flagCfg.Key, "k", "", "key")
	flag.Parse()

	log.Printf("Server init flags: %+v\n", flagCfg)

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	log.Printf("Server init ENV: %+v\n", cfg)

	if flagCfg.Address != address && cfg.Address == address {
		cfg.Address = flagCfg.Address
	}

	if flagCfg.Restore != restore && os.Getenv("RESTORE") == "" {
		cfg.Restore = flagCfg.Restore
	}

	if flagCfg.StoreInterval != storeInterval && cfg.StoreInterval == storeInterval {
		cfg.StoreInterval = flagCfg.StoreInterval
	}

	if flagCfg.StoreFile != storeFile && cfg.StoreFile == storeFile {
		cfg.StoreFile = flagCfg.StoreFile
	}

	if flagCfg.Key != key && cfg.Key == key {
		cfg.Key = flagCfg.Key
	}

	log.Printf("Server config: %+v\n", cfg)

	return cfg
}
