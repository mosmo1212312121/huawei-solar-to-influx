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
	MeterStatus = ModbusRegister{37100, 1, 0, "", "Meter Status"} // int16

	// --- Running Data ---
	// PV Input
	PV1Voltage = ModbusRegister{32016, 1, 10, "V", "PV1 Input Voltage"}  // int16
	PV1Current = ModbusRegister{32017, 1, 100, "A", "PV1 Input Current"} // int16
	PVPower    = ModbusRegister{32064, 2, 1, "W", "PV Power"}            // int32

	// Grid Output
	LineVoltageA     = ModbusRegister{37101, 2, 10, "V", "Line Voltage A"}    // int32
	LineVoltageB     = ModbusRegister{37103, 2, 10, "V", "Line Voltage B"}    // int32
	LineVoltageC     = ModbusRegister{37105, 2, 10, "V", "Line Voltage C"}    // int32
	PhaseACurrent    = ModbusRegister{37107, 2, 100, "A", "Phase A Current"}  // int32
	PhaseBCurrent    = ModbusRegister{37109, 2, 100, "A", "Phase B Current"}  // int32
	PhaseCCurrent    = ModbusRegister{37111, 2, 100, "A", "Phase C Current"}  // int32
	ActivePowerMeter = ModbusRegister{37113, 2, 1, "W", "Active Power meter"} // int32 ไฟที่ meter grid มาถ้า - คือจาก grid + คือย้อนออก
	PowerFactor      = ModbusRegister{37117, 1, 1000, "", "Power Factor"}     // int16
	GridFreq         = ModbusRegister{37118, 1, 100, "Hz", "Grid Frequency"}  // int16
	InverterPower    = ModbusRegister{32080, 2, 1, "W", "Inverter Power"}     // int32 ไฟที่ได้จาก inverter
	DeraRating       = ModbusRegister{40125, 1, 10, "%", "Dera Rating"}       // int16
)
