package sender

import (
	"github.com/Kuart/metric-collector/internal/metric"
	"net/http"
)

type Client struct {
	updatePath      string
	batchUpdatePath string
	pingPath        string
	isBatchEnable   bool
	crypto          EncryptionClient
	client          *http.Client
}

type EncryptionClient interface {
	EncodeMetric(m metric.Metric) metric.Metric
	Compare(h string, m metric.Metric) bool
}
