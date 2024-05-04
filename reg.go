package cpu8086

type Reg8 struct {
	rawMem    [RAM_SIZE]byte
	baseIndex int
}

type Reg16 struct {
	rawMem    [RAM_SIZE]byte
	baseIndex int
}

func NewReg8(baseIndex int, rawMem [RAM_SIZE]byte) *Reg8 {
	return &Reg8{
		rawMem:    rawMem,
		baseIndex: baseIndex,
	}
}

func NewReg16(baseIndex int, rawMem [RAM_SIZE]byte) *Reg16 {
	return &Reg16{
		rawMem:    rawMem,
		baseIndex: baseIndex,
	}
}

func (r *Reg8) Get(index int) byte {
	return r.rawMem[r.baseIndex+index]
}

func (r *Reg8) Set(index int, value byte) {
	r.rawMem[r.baseIndex+index] = value
}

func (r *Reg16) Get(index int) uint16 {
	return uint16(r.rawMem[r.baseIndex+index]) | uint16(r.rawMem[r.baseIndex+index+1])<<8
}

func (r *Reg16) Set(index int, value uint16) {
	r.rawMem[r.baseIndex+index] = byte(value)
	r.rawMem[r.baseIndex+index+1] = byte(value >> 8)
}
