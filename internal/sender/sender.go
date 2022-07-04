package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	config "github.com/Kuart/metric-collector/config/agent"
	"github.com/Kuart/metric-collector/internal/encryption"
	"github.com/Kuart/metric-collector/internal/metric"
	"log"
	"net/http"
)

func NewMetricClient(config config.Config, crypto encryption.Encryption) *Client {
	updatePath := fmt.Sprintf("http://%s/update", config.Address)
	batchUpdatePath := fmt.Sprintf("http://%s/updates", config.Address)
	return &Client{
		updatePath:      updatePath,
		batchUpdatePath: batchUpdatePath,
		crypto:          crypto,
		client: &http.Client{
			Timeout: config.PollInterval,
		},
	}
}

func (c *Client) SendMetrics(gauge metric.GaugeState, counter metric.Counter) {
	body := make([]metric.Metric, 0, len(gauge)+1)

	for key, value := range gauge {
		m := metric.NewMetricToSend(key, metric.GaugeTypeName, value)
		body = append(body, c.crypto.EncodeMetric(m))
	}

	body = append(body,
		c.crypto.EncodeMetric(metric.NewMetricToSend(counter.Name, metric.CounterTypeName, counter.Value)))

	if len(body) == 0 {
		log.Printf("nothing to send, metrics are empty")
		return
	}

	jsonValue, err := json.Marshal(body)

	if err != nil {
		log.Printf("metrics not sended, json marshal err: %s", err)
		return
	}

	buff := bytes.NewBuffer(jsonValue)
	response, err := c.client.Post(c.batchUpdatePath, "application/json;charset=utf-8", buff)

	if err != nil {
		log.Printf("metrics not sended, err: %s", err)
		return
	}

	response.Header.Set("Content-Encoding", "gzip")

	log.Printf("metrics batch sent, path %s", c.batchUpdatePath)

	defer response.Body.Close()
}

func (c *Client) sendGauge(gauge metric.GaugeState) {
	for key, value := range gauge {
		body := metric.NewMetricToSend(key, metric.GaugeTypeName, value)
		c.doRequest(body)
	}
}

func (c *Client) sendCounter(counter metric.Counter) {
	body := metric.NewMetricToSend(counter.Name, metric.CounterTypeName, counter.Value)
	c.doRequest(body)
}

func (c *Client) doRequest(body metric.Metric) {
	m := c.crypto.EncodeMetric(body)
	jsonValue, err := json.Marshal(m)

	if err != nil {
		log.Printf("%s metric not sended, json marshal err: %s", body.ID, err)
		return
	}

	buff := bytes.NewBuffer(jsonValue)
	response, err := c.client.Post(c.updatePath, "application/json;charset=utf-8", buff)

	if err != nil {
		log.Printf("%s metric not sended, err: %s", body.ID, err)
		return
	}

	log.Printf("%s metric sended", body.ID)

	defer response.Body.Close()
}
