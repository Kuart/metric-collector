package sender

import "net/http"

type Client struct {
	host       string
	port       string
	updatePath string
	client     *http.Client
}
