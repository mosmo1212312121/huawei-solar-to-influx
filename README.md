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

| Field              | Register  | Unit | Description                                  |
| ------------------ | --------- | ---- | -------------------------------------------- |
| PV1 Input Voltage  | 32016     | V    | PV string 1 voltage                          |
| PV1 Input Current  | 32017     | A    | PV string 1 current                          |
| PV Power           | 32064     | W    | Total PV input power                         |
| Inverter Power     | 32080     | W    | AC output power from inverter                |
| Line Voltage A     | 37101     | V    | Grid voltage phase A                         |
| Phase A Current    | 37107     | A    | Grid current phase A                         |
| Active Power meter | 37113     | W    | Grid meter power (− = from grid, + = export) |
| Power Factor       | 37117     | —    | Grid power factor                            |
| Grid Frequency     | 37118     | Hz   | Grid frequency                               |
| Dera Rating        | 40125     | %    | Active power derating percentage             |
| **Load Power**     | _derived_ | W    | `Inverter Power − Active Power meter`        |

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

1. เปิด Grafana → **Dashboards → Import**
2. Upload ไฟล์ `grafana_dashboard/solar_monitor.json`
3. เลือก InfluxDB datasource → **Import**

---

## Configuration

ค่า config ทั้งหมดอ่านจาก `.env` หรือ environment variable โดยตรง

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
