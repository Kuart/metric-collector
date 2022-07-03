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
	return &Client{
		updatePath: updatePath,
		crypto:     crypto,
		client: &http.Client{
			Timeout: config.PollInterval,
		},
	}
}

func (c *Client) SendMetrics(gauge metric.GaugeState, counter metric.Counter) {
	c.sendGauge(gauge)
	c.sendCounter(counter)
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
