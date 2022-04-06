package converter

import (
	"encoding/binary"
	"math"
)

func ConvertInt8ToBytes(i int8) []byte {
	var ret = make([]byte, 2)
	binary.BigEndian.PutUint16(ret, uint16(i))
	return ret
}

func ConvertInt32ToBytes(i int32) []byte {
	var ret = make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(i))
	return ret
}

func ConvertInt64ToBytes(i int64) []byte {
	var ret = make([]byte, 8)
	binary.BigEndian.PutUint64(ret, uint64(i))
	return ret
}

func ConvertFloat32ToBytes(f float32) []byte {
	var ret = make([]byte, 4)
	binary.BigEndian.PutUint32(ret, math.Float32bits(f))
	return ret
}

func ConvertFloat64ToBytes(f float64) []byte {
	var ret = make([]byte, 8)
	binary.BigEndian.PutUint64(ret, math.Float64bits(f))
	return ret
}

func ParseBytesToInt32(bs []byte) int32 {
	return int32(binary.BigEndian.Uint32(bs))
}

func ParseBytesToInt64(bs []byte) int64 {
	return int64(binary.BigEndian.Uint64(bs))
}

func ParseBytesToFloat32(bs []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(bs))
}

func ParseBytesToFloat64(bs []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(bs))
}
