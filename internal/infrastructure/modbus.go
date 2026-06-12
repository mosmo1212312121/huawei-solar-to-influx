package infrastructure

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/goburrow/modbus"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/config"
	"github.com/mosmo1212312121/huawei-solar-to-influx/internal/utils"
)

type ModbusClient struct {
	cfg     config.ModbusConfig
	handler *modbus.TCPClientHandler
	Client  modbus.Client
}

func ConnectModbus(cfg config.ModbusConfig) (*ModbusClient, error) {
	handler := modbus.NewTCPClientHandler(cfg.Addr)
	handler.Timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	if err := handler.Connect(); err != nil {
		return nil, err
	}
	log.Println("Modbus connected successfully")
	return &ModbusClient{
		cfg:     cfg,
		handler: handler,
		Client:  modbus.NewClient(handler),
	}, nil
}

func (m *ModbusClient) Reconnect() {
	log.Println("Reconnecting Modbus...")
	_ = m.handler.Close()
	time.Sleep(time.Duration(m.cfg.ReconnectWaitSeconds) * time.Second)

	newHandler := modbus.NewTCPClientHandler(m.cfg.Addr)
	newHandler.Timeout = time.Duration(m.cfg.TimeoutSeconds) * time.Second
	if err := newHandler.Connect(); err != nil {
		log.Printf("Modbus reconnect error: %v", err)
		return
	}
	m.handler = newHandler
	m.Client = modbus.NewClient(newHandler)
	log.Println("Modbus reconnected successfully")
}

func (m *ModbusClient) Close() {
	_ = m.handler.Close()
}

// IsConnError reports whether err is a transport/connection failure (timeout,
// EOF, reset) rather than a Modbus protocol exception. A *modbus.ModbusError
// means the device answered but rejected the request — e.g. the register does
// not exist on this model — so the connection is fine and we should not reconnect.
func IsConnError(err error) bool {
	if err == nil {
		return false
	}
	var modbusErr *modbus.ModbusError
	return !errors.As(err, &modbusErr)
}

func ReadRegister(c modbus.Client, reg utils.ModbusRegister) (float64, error) {
	raw, err := c.ReadHoldingRegisters(reg.Address, reg.Quantity)
	if err != nil {
		return 0, err
	}

	var intVal int64
	switch reg.Quantity {
	case 2:
		intVal = int64(utils.ToInt32(raw))
	default: // 1 register → int16
		intVal = int64(utils.ToInt16(raw))
	}

	// 0x7FFFFFFF / 0x7FFF = Huawei sentinel for "data not available"
	if intVal == math.MaxInt32 || intVal == math.MaxInt16 {
		return 0, nil
	}

	if reg.Gain == 0 {
		return float64(intVal), nil
	}
	return float64(intVal) / float64(reg.Gain), nil
}
