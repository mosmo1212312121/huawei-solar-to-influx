package utils

import (
	"encoding/binary"
	"fmt"
)

func ToUint16(b []byte) uint16 {
	if len(b) < 2 {
		return 0
	}
	return binary.BigEndian.Uint16(b[:2])
}

func ToInt16(b []byte) int16 {
	return int16(ToUint16(b))
}
func ToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return binary.BigEndian.Uint32(b[:4])
}

func ToInt32(b []byte) int32 {
	return int32(ToUint32(b))
}
func ToString(b []byte) string {
	return fmt.Sprintf("%s", b)
}

func ToFloat32WithGain(b []byte, gain float64) float64 {
	raw := ToInt16(b)
	return float64(raw) / gain
}
