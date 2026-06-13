# huawei-solar-to-influx

A Go service that polls a **Huawei SUN2000 inverter** via **Modbus TCP**, computes derived metrics, and writes the data to **InfluxDB 1.x** for visualization in **Grafana**.

---

## Architecture

```
Huawei SUN2000
  (Modbus TCP)
       │
       ▼
huawei-solar-to-influx  ──►  InfluxDB 1.8  ──►  Grafana
```

---

## Metrics Collected


Register definitions [`internal/utils/register.go`](internal/utils/register.go).

### State / Alarm

| Field           | Register | Unit | Description                                          |
| --------------- | -------- | ---- | ---------------------------------------------------- |
| Meter Status    | 37100    | —    | Power meter status                                   |
| Device Status   | 32089    | —    | Inverter status (0x0000 standby, 0x0200 on-grid, …)  |
| Fault Code      | 32090    | —    | Inverter fault code                                  |
| Inverter Alarm 1| 32008    | —    | Alarm bitfield 1                                     |
| Inverter Alarm 2| 32009    | —    | Alarm bitfield 2                                     |
| Inverter Alarm 3| 32010    | —    | Alarm bitfield 3                                     |

### PV Input

| Field             | Register | Unit | Description           |
| ----------------- | -------- | ---- | --------------------- |
| PV1 Input Voltage | 32016    | V    | PV string 1 voltage   |
| PV1 Input Current | 32017    | A    | PV string 1 current   |
| PV2 Input Voltage | 32018    | V    | PV string 2 voltage   |
| PV2 Input Current | 32019    | A    | PV string 2 current   |
| PV3 Input Voltage | 32020    | V    | PV string 3 voltage   |
| PV3 Input Current | 32021    | A    | PV string 3 current   |
| PV4 Input Voltage | 32022    | V    | PV string 4 voltage   |
| PV4 Input Current | 32023    | A    | PV string 4 current   |
| PV Power          | 32064    | W    | Total PV input power  |

### Inverter Output

| Field                    | Register | Unit | Description                   |
| ------------------------ | -------- | ---- | ----------------------------- |
| Inverter Phase A Current | 32072    | A    | Phase A current               |
| Inverter Phase B Current | 32074    | A    | Phase B current               |
| Inverter Phase C Current | 32076    | A    | Phase C current               |
| Peak Active Power of Day | 32078    | W    | Peak active power today       |
| Inverter Power           | 32080    | W    | AC output power from inverter |
| Inverter Reactive Power  | 32082    | kvar | Reactive power                |
| Inverter Power Factor    | 32084    | —    | Inverter power factor         |
| Inverter Frequency       | 32085    | Hz   | Inverter output frequency     |
| Inverter Efficiency      | 32086    | %    | Conversion efficiency         |
| Internal Temperature     | 32087    | °C   | Internal temperature          |
| Insulation Resistance    | 32088    | MΩ   | Insulation resistance         |

### Energy Yield

| Field                    | Register | Unit | Description                  |
| ------------------------ | -------- | ---- | ---------------------------- |
| Accumulated Energy Yield | 32106    | kWh  | Lifetime energy produced     |
| Daily Energy Yield       | 32114    | kWh  | Energy produced today        |

### Grid (Power Meter)

| Field                    | Register | Unit | Description                                  |
| ------------------------ | -------- | ---- | -------------------------------------------- |
| Line Voltage A           | 37101    | V    | Grid voltage phase A                         |
| Line Voltage B           | 37103    | V    | Grid voltage phase B                         |
| Line Voltage C           | 37105    | V    | Grid voltage phase C                         |
| Phase A Current          | 37107    | A    | Grid current phase A                         |
| Phase B Current          | 37109    | A    | Grid current phase B                         |
| Phase C Current          | 37111    | A    | Grid current phase C                         |
| Active Power meter       | 37113    | W    | Grid meter power (− = from grid, + = export) |
| Reactive Power meter     | 37115    | var  | Grid reactive power                          |
| Power Factor             | 37117    | —    | Grid power factor                            |
| Grid Frequency           | 37118    | Hz   | Grid frequency                               |
| Grid Exported Energy     | 37119    | kWh  | Cumulative energy exported to grid           |
| Grid Imported Energy     | 37121    | kWh  | Cumulative energy imported from grid         |

### Battery / Energy Storage (LUNA2000)

| Field                          | Register | Unit | Description                                     |
| ------------------------------ | -------- | ---- | ----------------------------------------------- |
| Battery SOC                    | 37760    | %    | State of charge                                 |
| Battery Running Status         | 37762    | —    | 0 offline · 1 standby · 2 running · 3 fault · 4 sleep |
| Battery Charge/Discharge Power | 37765    | W    | + = charging, − = discharging                   |
| Battery Total Charge           | 37780    | kWh  | Cumulative energy charged                       |
| Battery Total Discharge        | 37782    | kWh  | Cumulative energy discharged                    |

### Settings & Derived

| Field          | Register  | Unit | Description                           |
| -------------- | --------- | ---- | ------------------------------------- |
| Dera Rating    | 40125     | %    | Active power derating percentage      |
| **Load Power** | _derived_ | W    | `Inverter Power − Active Power meter` |

---

## Quick Start

### 1. Clone & configure

```bash
git clone https://github.com/mosmo1212312121/huawei-solar-to-influx.git
cd huawei-solar-to-influx
cp .env.example .env
```

### 2. Run with Docker Compose

```bash
docker compose up -d
```

Services :

| Service                | URL                   |
| ---------------------- | --------------------- |
| InfluxDB               | http://localhost:8086 |
| Grafana                | http://localhost:3000 |
| huawei-solar-to-influx | — (background poller) |

### 3. Import Grafana Dashboard

1. Open Grafana → **Dashboards → Import**
2. Upload the `grafana_dashboard/solar_monitor.json` file
3. Select the InfluxDB datasource → **Import**

---

## Configuration

All configuration is read directly from `.env` or environment variables.

| Variable                        | Default                   | Description            |
| ------------------------------- | ------------------------- | ---------------------- |
| `MODBUS_ADDR`                   | `192.168.200.1:6607`      | Inverter IP:port       |
| `MODBUS_TIMEOUT_SECONDS`        | `5`                       | Modbus read timeout    |
| `MODBUS_RECONNECT_WAIT_SECONDS` | `2`                       | Wait before reconnect  |
| `INFLUX_URL`                    | `http://192.168.1.5:8086` | InfluxDB URL           |
| `INFLUX_DB`                     | `solar`                   | InfluxDB database name |
| `INFLUX_MEASUREMENT`            | `huawei_solar`            | Measurement name       |
| `INFLUX_DEVICE_TAG`             | `SUN2000`                 | `device` tag value     |
| `APP_POLL_INTERVAL_SECONDS`     | `10`                      | Polling interval       |

---

## Project Structure

```
.
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Viper config loader
│   ├── infrastructure/
│   │   ├── modbus.go            # Modbus connect / reconnect / read
│   │   └── influx.go            # InfluxDB connect
│   └── utils/
│       ├── register.go          # Modbus register definitions
│       └── converter_base.go    # int16 / int32 converters
├── grafana_dashboard/
│   └── solar_monitor.json       # Grafana dashboard export
├── .env.example                 # Config template
├── dockerfile
├── docker-compose.yml
└── .github/
    └── workflows/
        └── docker-publish.yml   # CI/CD → ghcr.io
```

---

## Docker Image

Pre-built images are published to GitHub Container Registry on every version tag.

```bash
docker pull ghcr.io/mosmo1212312121/huawei-solar-to-influx:latest
```

**Supported platforms:** `linux/amd64` · `linux/arm64` · `linux/arm/v7`

---

## Release

```bash
git tag v1.0.0
git push origin v1.0.0
```
