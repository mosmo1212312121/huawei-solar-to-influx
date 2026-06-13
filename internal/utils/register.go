package utils

type ModbusRegister struct {
	Address  uint16
	Quantity uint16
	Gain     int16
	Unit     string
	Desc     string
}

var (
	// --- State ---
	MeterStatus    = ModbusRegister{37100, 1, 0, "", "Meter Status"}     // int16
	DeviceStatus   = ModbusRegister{32089, 1, 0, "", "Device Status"}    // uint16 (0x0000 standby, 0x0200 on-grid, 0x0300 shutdown ฯลฯ)
	FaultCode      = ModbusRegister{32090, 1, 0, "", "Fault Code"}       // uint16
	InverterAlarm1 = ModbusRegister{32008, 1, 0, "", "Inverter Alarm 1"} // uint16 bitfield
	InverterAlarm2 = ModbusRegister{32009, 1, 0, "", "Inverter Alarm 2"} // uint16 bitfield
	InverterAlarm3 = ModbusRegister{32010, 1, 0, "", "Inverter Alarm 3"} // uint16 bitfield

	// --- Running Data ---
	// PV Input
	PV1Voltage = ModbusRegister{32016, 1, 10, "V", "PV1 Input Voltage"}  // int16
	PV1Current = ModbusRegister{32017, 1, 100, "A", "PV1 Input Current"} // int16
	PV2Voltage = ModbusRegister{32018, 1, 10, "V", "PV2 Input Voltage"}  // int16
	PV2Current = ModbusRegister{32019, 1, 100, "A", "PV2 Input Current"} // int16
	PV3Voltage = ModbusRegister{32020, 1, 10, "V", "PV3 Input Voltage"}  // int16
	PV3Current = ModbusRegister{32021, 1, 100, "A", "PV3 Input Current"} // int16
	PV4Voltage = ModbusRegister{32022, 1, 10, "V", "PV4 Input Voltage"}  // int16
	PV4Current = ModbusRegister{32023, 1, 100, "A", "PV4 Input Current"} // int16
	PVPower    = ModbusRegister{32064, 2, 1, "W", "PV Power"}            // int32

	// Inverter Output
	InverterCurrentA     = ModbusRegister{32072, 2, 1000, "A", "Inverter Phase A Current"}   // int32
	InverterCurrentB     = ModbusRegister{32074, 2, 1000, "A", "Inverter Phase B Current"}   // int32
	InverterCurrentC     = ModbusRegister{32076, 2, 1000, "A", "Inverter Phase C Current"}   // int32
	PeakActivePowerDay   = ModbusRegister{32078, 2, 1, "W", "Peak Active Power of Day"}      // int32
	InverterPower        = ModbusRegister{32080, 2, 1, "W", "Inverter Power"}                // int32 ไฟที่ได้จาก inverter
	InverterReactivePwr  = ModbusRegister{32082, 2, 1000, "kvar", "Inverter Reactive Power"} // int32
	InverterPowerFactor  = ModbusRegister{32084, 1, 1000, "", "Inverter Power Factor"}       // int16
	InverterFreq         = ModbusRegister{32085, 1, 100, "Hz", "Inverter Frequency"}         // uint16
	InverterEfficiency   = ModbusRegister{32086, 1, 100, "%", "Inverter Efficiency"}         // uint16
	InternalTemperature  = ModbusRegister{32087, 1, 10, "C", "Internal Temperature"}         // int16
	InsulationResistance = ModbusRegister{32088, 1, 1000, "MOhm", "Insulation Resistance"}   // uint16

	// Energy Yield
	AccumulatedEnergyYield = ModbusRegister{32106, 2, 100, "kWh", "Accumulated Energy Yield"} // uint32 พลังงานสะสมทั้งหมด
	DailyEnergyYield       = ModbusRegister{32114, 2, 100, "kWh", "Daily Energy Yield"}       // uint32 พลังงานที่ผลิตได้วันนี้

	// Grid (Power Meter)
	LineVoltageA       = ModbusRegister{37101, 2, 10, "V", "Line Voltage A"}          // int32
	LineVoltageB       = ModbusRegister{37103, 2, 10, "V", "Line Voltage B"}          // int32
	LineVoltageC       = ModbusRegister{37105, 2, 10, "V", "Line Voltage C"}          // int32
	PhaseACurrent      = ModbusRegister{37107, 2, 100, "A", "Phase A Current"}        // int32
	PhaseBCurrent      = ModbusRegister{37109, 2, 100, "A", "Phase B Current"}        // int32
	PhaseCCurrent      = ModbusRegister{37111, 2, 100, "A", "Phase C Current"}        // int32
	ActivePowerMeter   = ModbusRegister{37113, 2, 1, "W", "Active Power meter"}       // int32 ไฟที่ meter grid มาถ้า - คือจาก grid + คือย้อนออก
	ReactivePowerMeter = ModbusRegister{37115, 2, 1, "var", "Reactive Power meter"}   // int32
	PowerFactor        = ModbusRegister{37117, 1, 1000, "", "Power Factor"}           // int16
	GridFreq           = ModbusRegister{37118, 1, 100, "Hz", "Grid Frequency"}        // int16
	GridExportedEnergy = ModbusRegister{37119, 2, 100, "kWh", "Grid Exported Energy"} // int32 พลังงานที่ขายออก grid สะสม
	GridImportedEnergy = ModbusRegister{37121, 2, 100, "kWh", "Grid Imported Energy"} // int32 พลังงานที่ซื้อจาก grid สะสม

	// Battery / Energy Storage (LUNA2000 — ใช้ได้เฉพาะรุ่นที่ต่อแบตเตอรี่)
	BatteryRunningStatus  = ModbusRegister{37762, 1, 0, "", "Battery Running Status"}          // uint16 (0 offline, 1 standby, 2 running, 3 fault, 4 sleep)
	BatteryChargePower    = ModbusRegister{37765, 2, 1, "W", "Battery Charge/Discharge Power"} // int32 + ชาร์จ, - คาย
	BatterySOC            = ModbusRegister{37760, 1, 10, "%", "Battery SOC"}                   // uint16
	BatteryTotalCharge    = ModbusRegister{37780, 2, 100, "kWh", "Battery Total Charge"}       // uint32
	BatteryTotalDischarge = ModbusRegister{37782, 2, 100, "kWh", "Battery Total Discharge"}    // uint32

	// Settings
	DeraRating = ModbusRegister{40125, 1, 10, "%", "Dera Rating"} // int16
)
