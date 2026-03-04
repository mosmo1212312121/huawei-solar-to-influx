package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ModbusConfig struct {
	Addr                 string `mapstructure:"addr"`
	TimeoutSeconds       int    `mapstructure:"timeout_seconds"`
	ReconnectWaitSeconds int    `mapstructure:"reconnect_wait_seconds"`
}

type InfluxConfig struct {
	URL         string `mapstructure:"url"`
	DB          string `mapstructure:"db"`
	DeviceTag   string `mapstructure:"device_tag"`
	Measurement string `mapstructure:"measurement"`
}

type AppConfig struct {
	PollIntervalSeconds int `mapstructure:"poll_interval_seconds"`
}

type Config struct {
	Modbus ModbusConfig `mapstructure:"modbus"`
	Influx InfluxConfig `mapstructure:"influx"`
	App    AppConfig    `mapstructure:"app"`
}

func Load() (*Config, error) {
	// Load .env file if present (silently ignored if not found)
	_ = godotenv.Load()

	// Map MODBUS_ADDR → modbus.addr, INFLUX_URL → influx.url, etc.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults (used when env var is not set)
	viper.SetDefault("modbus.addr", "192.168.200.1:6607")
	viper.SetDefault("modbus.timeout_seconds", 5)
	viper.SetDefault("modbus.reconnect_wait_seconds", 2)
	viper.SetDefault("influx.url", "http://192.168.1.5:8086")
	viper.SetDefault("influx.db", "solar")
	viper.SetDefault("influx.device_tag", "SUN2000")
	viper.SetDefault("influx.measurement", "huawei_solar")
	viper.SetDefault("app.poll_interval_seconds", 10)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
