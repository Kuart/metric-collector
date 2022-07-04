package sender

import (
	"github.com/Kuart/metric-collector/internal/encryption"
	"net/http"
)

type Client struct {
	updatePath      string
	batchUpdatePath string
	crypto          encryption.Encryption
	client          *http.Client
}
