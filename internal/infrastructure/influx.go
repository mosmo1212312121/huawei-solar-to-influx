package infrastructure

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/config"
)

func ConnectInflux(cfg config.InfluxConfig) (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: cfg.URL,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
