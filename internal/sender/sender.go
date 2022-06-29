package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kuart/metric-collector/internal/metric"
	"log"
	"net/http"
	"time"
)

func NewMetricClient(address string, pollInterval time.Duration) *Client {
	updatePath := fmt.Sprintf("http://%s/update", address)
	return &Client{
		updatePath: updatePath,
		client: &http.Client{
			Timeout: pollInterval,
		},
	}
}

func (c *Client) SendMetrics(gauge metric.GaugeState, counter metric.Counter) {
	c.sendGauge(gauge)
	c.sendCounter(counter)
}

func (c *Client) sendGauge(gauge metric.GaugeState) {
	for key, value := range gauge {
		body := metric.Metric{
			ID:    key,
			MType: metric.GaugeTypeName,
			Value: &value,
		}

		c.doRequest(body)
	}
}

func (c *Client) sendCounter(counter metric.Counter) {
	body := metric.Metric{
		ID:    counter.Name,
		MType: metric.CounterTypeName,
		Delta: &counter.Value,
	}

	c.doRequest(body)
}

func (c *Client) doRequest(body metric.Metric) {
	jsonValue, err := json.Marshal(body)

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
