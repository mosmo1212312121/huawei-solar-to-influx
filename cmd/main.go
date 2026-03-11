package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/config"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/infrastructure"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config load error:", err)
	}

	// --- InfluxDB ---
	ic, err := infrastructure.ConnectInflux(cfg.Influx)
	if err != nil {
		log.Fatal("InfluxDB connect error:", err)
	}
	defer ic.Close()

	// --- Modbus TCP ---
	mbClient, err := infrastructure.ConnectModbus(cfg.Modbus)
	if err != nil {
		log.Fatal("Modbus connect error:", err)
	}
	defer mbClient.Close()

	// --- Registers to poll ---
	registers := []utils.ModbusRegister{
		utils.MeterStatus,
		utils.PV1Voltage,
		utils.PV1Current,
		utils.PVPower,
		utils.LineVoltageA,
		utils.PhaseACurrent,
		utils.ActivePowerMeter,
		utils.PowerFactor,
		utils.GridFreq,
		utils.InverterPower,
		utils.DeraRating,
	}

	pollInterval := time.Duration(cfg.App.PollIntervalSeconds) * time.Second

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	log.Printf("Polling every %v ...", pollInterval)

	for {
		select {
		case <-quit:
			log.Println("Shutting down...")
			mbClient.Close()
			log.Println("modbus client closed")
			ic.Close()
			log.Println("influxdb client closed")
			return

		case <-ticker.C:
			fields := make(map[string]interface{})
			hasError := false

			for _, reg := range registers {
				val, err := infrastructure.ReadRegister(mbClient.Client, reg)
				if err != nil {
					log.Printf("Modbus read error [%s addr=%d]: %v", reg.Desc, reg.Address, err)
					hasError = true
					mbClient.Reconnect()
					break
				}
				fields[reg.Desc] = val
				log.Printf("  %-25s = %.4f %s", reg.Desc, val, reg.Unit)
			}

			if hasError {
				log.Println("Read error occurred — skipping write, will retry next poll")
				continue
			}
			if len(fields) == 0 {
				log.Println("No data collected, skipping write")
				continue
			}

			// --- Write to InfluxDB ---
			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  cfg.Influx.DB,
				Precision: "s",
			})
			if err != nil {
				log.Printf("BatchPoints error: %v", err)
				continue
			}

			// Calculate Load Power
			inverterPower := fields[utils.InverterPower.Desc].(float64)
			activePowerMeter := fields[utils.ActivePowerMeter.Desc].(float64)
			loadPower := calcLoadPower(inverterPower, activePowerMeter)
			fields["Load Power"] = loadPower

			// Calculate Dera Rating
			// deraRatingPercentage := fields[utils.DeraRating.Desc].(float64)
			// deraRatingWatts := deraRatingToWatts(deraRatingPercentage, inverterPower)
			// fields["DeraRatingWatts"] = deraRatingWatts

			tags := map[string]string{"device": cfg.Influx.DeviceTag}
			pt, err := client.NewPoint(cfg.Influx.Measurement, tags, fields, time.Now())
			if err != nil {
				log.Printf("NewPoint error: %v", err)
				continue
			}
			bp.AddPoint(pt)

			if err := ic.Write(bp); err != nil {
				log.Printf("InfluxDB write error: %v", err)
			} else {
				log.Printf("Wrote %d field(s) to InfluxDB", len(fields))
			}
		}
	}
}
func calcLoadPower(inverterPower, activePowerMeter float64) float64 {
	loadPower := inverterPower - activePowerMeter
	log.Printf("  %-25s = %.4f %s", "Load Power", loadPower, "W")
	return loadPower
}

// Comment out for now, derating is not correct
// func deraRatingToWatts(deraRatingPercentage float64, inverterPower float64) float64 {
// 	deraRatingWatts := inverterPower * deraRatingPercentage / 100
// 	log.Printf("  %-25s = %.4f %s", "Dera Rating", deraRatingWatts, "W")
// 	return deraRatingWatts
// }
