package loki

import (
	"time"

	"github.com/go-kit/log"

	"github.com/grafana/dskit/flagext"
	"github.com/grafana/loki/clients/pkg/promtail/api"
	loki_client "github.com/grafana/loki/clients/pkg/promtail/client"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

type Client interface {
	PostLog(text string, timestamp time.Time, labels model.LabelSet)
}

type ClientParams struct {
	URI    string
	Logger log.Logger
}

type lclient struct {
	lokiClient loki_client.Client
}

func NewClient(params ClientParams) (Client, error) {

	url := flagext.URLValue{}
	err := url.Set(params.URI)
	if err != nil {
		return nil, err
	}

	client_config := loki_client.Config{
		URL:     url,
		Timeout: time.Minute,
	}

	clientMetrics := loki_client.NewMetrics(prometheus.DefaultRegisterer, client_config.StreamLagLabels)
	lokiClient, err := loki_client.NewMulti(clientMetrics, client_config.StreamLagLabels, params.Logger, client_config)
	if err != nil {
		return nil, err
	}

	return &lclient{lokiClient}, nil
}

func (c *lclient) PostLog(text string, timestamp time.Time, labels model.LabelSet) {
	c.lokiClient.Chan() <- api.Entry{
		Labels: labels,
		Entry: logproto.Entry{
			Line:      text,
			Timestamp: timestamp,
		},
	}
}
