package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Kuart/metric-collector/internal/metric"
)

type Encryption struct {
	key   []byte
	IsKey bool
}

func New(k string) Encryption {
	return Encryption{
		key:   []byte(k),
		IsKey: k != "",
	}
}

func (c Encryption) EncodeMetric(m metric.Metric) metric.Metric {
	if !c.IsKey {
		return m
	}

	var src string

	if m.MType == metric.GaugeTypeName {
		src = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	} else if m.MType == metric.CounterTypeName {
		src = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}

	h := hmac.New(sha256.New, c.key)
	h.Write([]byte(src))
	m.Hash = hex.EncodeToString(h.Sum(nil))

	return m
}

func (c Encryption) Compare(h string, m metric.Metric) bool {
	return hmac.Equal([]byte(h), []byte(c.EncodeMetric(m).Hash))
}
