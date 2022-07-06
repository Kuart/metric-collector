package sender

import (
	"github.com/Kuart/metric-collector/internal/encryption"
	"net/http"
)

type Client struct {
	updatePath      string
	batchUpdatePath string
	pingPath        string
	isBatchEnable   bool
	crypto          encryption.Encryption
	client          *http.Client
}
