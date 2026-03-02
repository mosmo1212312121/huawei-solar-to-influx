package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/infrastructure"

	"github.com/goburrow/modbus"
)

func main() {

	// InfluxDB
	dbName := "my_database"
	c, err := infrastructure.ConnectInflux()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// --- Modbus TCP ---
	handler := modbus.NewTCPClientHandler("127.0.0.1:10502")
	handler.Timeout = 5 * time.Second
	if err := handler.Connect(); err != nil {
		log.Fatal("Modbus Error:", err)
	}
	defer handler.Close()
	mbClient := modbus.NewClient(handler)

	//

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		results, err := mbClient.ReadHoldingRegisters(0, 1)
		if err != nil {
			log.Printf("Modbus Read Error: %v", err)
			continue
		}

		// แปลงค่า (สมมติว่าเป็น Integer 16-bit)
		val := int16(results[0])<<8 | int16(results[1])

		// --- 4. เขียนข้อมูลลง InfluxDB 1.8 ---
		// สร้าง Batch Points
		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  dbName,
			Precision: "s", // วินาที
		})

		// สร้าง Point
		tags := map[string]string{"device": "PLC01"}
		fields := map[string]interface{}{"temperature": float64(val)}
		pt, _ := client.NewPoint("sensor_data", tags, fields, time.Now())

		bp.AddPoint(pt)

		// เขียนข้อมูล
		if err := c.Write(bp); err != nil {
			log.Printf("InfluxDB Write Error: %v", err)
		} else {
			log.Printf("Data sent: %v", val)
		}
	}
	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	log.Println("Server exited")
}
