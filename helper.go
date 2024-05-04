package cpu8086

type hasZero interface {
	byte |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Bool[T hasZero](v T) T {
	if v == 0 {
		return 0
	}
	return 1
}

func BoolByte[T hasZero](v T) byte {
	if v == 0 {
		return 0
	}
	return 1
}
