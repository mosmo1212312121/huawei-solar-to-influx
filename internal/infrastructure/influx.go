package infrastructure

import (
	client "github.com/influxdata/influxdb1-client/v2"
)

func ConnectInflux() (client.Client, error) {
	influxURL := "http://localhost:8086"
	// dbName := "my_database"

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxURL,
	})
	if err != nil {
		return nil, err
	}

	return c, nil

	// bp, err := client.NewBatchPoints(client.BatchPointsConfig{
	// 	Database: dbName,
	// })
	// if err != nil {
	// 	return err
	// }

	// tags := map[string]string{
	// 	"host": "server01",
	// }

	// fields := map[string]interface{}{
	// 	"temperature": 25.3,
	// 	"humidity":    60,
	// }

	// pt, err := client.NewPoint(
	// 	"weather",
	// 	tags,
	// 	fields,
	// 	time.Now(),
	// )
	// if err != nil {
	// 	return err
	// }

	// bp.AddPoint(pt)

	// err = c.Write(bp)
	// if err != nil {
	// 	return err
	// }

	// log.Println("Write success")

}
