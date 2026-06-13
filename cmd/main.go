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
		// State / Alarm
		utils.MeterStatus,
		utils.DeviceStatus,
		utils.FaultCode,
		utils.InverterAlarm1,
		utils.InverterAlarm2,
		utils.InverterAlarm3,

		// PV Input
		utils.PV1Voltage,
		utils.PV1Current,
		utils.PV2Voltage,
		utils.PV2Current,
		utils.PV3Voltage,
		utils.PV3Current,
		utils.PV4Voltage,
		utils.PV4Current,
		utils.PVPower,

		// Inverter Output
		utils.InverterCurrentA,
		utils.InverterCurrentB,
		utils.InverterCurrentC,
		utils.PeakActivePowerDay,
		utils.InverterPower,
		utils.InverterReactivePwr,
		utils.InverterPowerFactor,
		utils.InverterFreq,
		utils.InverterEfficiency,
		utils.InternalTemperature,
		utils.InsulationResistance,

		// Energy Yield
		utils.AccumulatedEnergyYield,
		utils.DailyEnergyYield,

		// Grid (Power Meter)
		utils.LineVoltageA,
		utils.LineVoltageB,
		utils.LineVoltageC,
		utils.PhaseACurrent,
		utils.PhaseBCurrent,
		utils.PhaseCCurrent,
		utils.ActivePowerMeter,
		utils.ReactivePowerMeter,
		utils.PowerFactor,
		utils.GridFreq,
		utils.GridExportedEnergy,
		utils.GridImportedEnergy,

		// Battery / Energy Storage (LUNA2000)
		utils.BatteryRunningStatus,
		utils.BatteryChargePower,
		utils.BatterySOC,
		utils.BatteryTotalCharge,
		utils.BatteryTotalDischarge,

		// Settings
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
			connLost := false

			for _, reg := range registers {
				val, err := infrastructure.ReadRegister(mbClient.Client, reg)
				if err != nil {
					// Skip this register. A register can be absent on this model
					// (e.g. battery registers without a LUNA2000) — don't abort the
					// whole round. Reconnect only if the connection itself dropped.
					log.Printf("Modbus read error [%s addr=%d]: %v", reg.Desc, reg.Address, err)
					if infrastructure.IsConnError(err) {
						connLost = true
						mbClient.Reconnect()
						break
					}
					continue
				}
				// Skip missing/zero values — don't store them in InfluxDB.
				if val == 0 {
					log.Printf("  %-25s = 0 %s (skipped)", reg.Desc, reg.Unit)
					continue
				}
				fields[reg.Desc] = val
				log.Printf("  %-25s = %.4f %s", reg.Desc, val, reg.Unit)
			}

			if connLost {
				log.Println("Connection lost — skipping write, will retry next poll")
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

			// Calculate Load Power (only if both inputs were read this round)
			inverterPower, okInv := fields[utils.InverterPower.Desc].(float64)
			activePowerMeter, okMeter := fields[utils.ActivePowerMeter.Desc].(float64)
			if okInv && okMeter {
				loadPower := calcLoadPower(inverterPower, activePowerMeter)
				fields["Load Power"] = loadPower
			}

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
