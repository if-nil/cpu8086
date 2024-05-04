package cpu8086

import "time"

const (
	// Emulator system constants
	IO_PORT_COUNT  = 0x10000
	RAM_SIZE       = 0x10FFF0
	REGS_BASE      = 0xF0000
	VIDEO_RAM_SIZE = 0x10000

	// Graphics/timer/keyboard update delays (explained later)
	KEYBOARD_TIMER_UPDATE_DELAY = 20000

	// 16-bit register decodes
	REG_AX = 0
	REG_CX = 1
	REG_DX = 2
	REG_BX = 3
	REG_SP = 4
	REG_BP = 5
	REG_SI = 6
	REG_DI = 7

	REG_ES = 8
	REG_CS = 9
	REG_SS = 10
	REG_DS = 11

	// 8-bit register decodes
	REG_AL = 0
	REG_AH = 1
	REG_CL = 2
	REG_CH = 3
	REG_DL = 4
	REG_DH = 5
	REG_BL = 6
	REG_BH = 7

	// FLAGS register decodes
	FLAG_CF = 40
	FLAG_PF = 41
	FLAG_AF = 42
	FLAG_ZF = 43
	FLAG_SF = 44
	FLAG_TF = 45
	FLAG_IF = 46
	FLAG_DF = 47
	FLAG_OF = 48

	// Lookup tables in the BIOS binary
	TABLE_XLAT_OPCODE        = 8
	TABLE_XLAT_SUBFUNCTION   = 9
	TABLE_STD_FLAGS          = 10
	TABLE_PARITY_FLAG        = 11
	TABLE_BASE_INST_SIZE     = 12
	TABLE_I_W_SIZE           = 13
	TABLE_I_MOD_SIZE         = 14
	TABLE_COND_JUMP_DECODE_A = 15
	TABLE_COND_JUMP_DECODE_B = 16
	TABLE_COND_JUMP_DECODE_C = 17
	TABLE_COND_JUMP_DECODE_D = 18
	TABLE_FLAGS_BITFIELDS    = 19

	// Bitfields for TABLE_STD_FLAGS values
	FLAGS_UPDATE_SZP      = 1
	FLAGS_UPDATE_AO_ARITH = 2
	FLAGS_UPDATE_OC_LOGIC = 4
)

type Cpu struct {
	mem             [RAM_SIZE]byte
	ioPorts         [IO_PORT_COUNT]byte
	opcodeStream    []byte
	regs8           []byte
	iRm             byte
	iW              byte
	iReg            byte
	iMod            byte
	iModSize        byte
	iD              byte
	iReg4Bit        byte
	rawOpcodeId     byte
	xlatOpcodeId    byte
	extra           byte
	repMode         byte
	segOverrideEn   byte
	repOverrideEn   byte
	trapFlag        byte
	int8Asap        byte
	scratchUchar    byte
	ioHiLo          byte
	vidMemBase      []byte
	spkrEn          byte
	biosTableLookup [20][256]byte

	regs16      []uint16
	regIp       uint16
	segOverride uint16
	fileIndex   uint16
	waveCounter uint16

	opSource     uint32
	opDest       uint32
	rmAddr       uint32
	opToAddr     uint32
	opFromAddr   uint32
	iData0       uint32
	iData1       uint32
	iData2       uint32
	scratchUint  uint32
	scratch2Uint uint32
	instCounter  uint32
	setFlagsType uint32
	graphicsX    uint32
	graphicsY    uint32
	pixelColors  [16]uint32
	vmemCtr      uint32

	opResult   int
	disk       [3]int
	scratchInt int

	clockBuf time.Time
	msClock  time.Time
}

func NewCPU() *Cpu {
	cpu := &Cpu{}
	cpu.regs8 = make([]byte, 0)
	cpu.regs16 = make([]uint16, 0)
	cpu.vidMemBase = make([]byte, 0)
	return cpu
}

// Returns number of top bit in operand (i.e. 8 for 8-bit operands, 16 for 16-bit operands)
func (c *Cpu) topBit() byte {
	return 8 * (c.iW + 1)
}

// Set carry flag
func (c *Cpu) setCf(newCf int) byte {
	if newCf != 0 {
		c.regs8[FLAG_CF] = 1
	} else {
		c.regs8[FLAG_CF] = 0
	}
	return c.regs8[FLAG_CF]
}

// Set auxiliary flag
func (c *Cpu) setAf(newAf int) byte {
	if newAf != 0 {
		c.regs8[FLAG_AF] = 1
	} else {
		c.regs8[FLAG_AF] = 0
	}
	return c.regs8[FLAG_AF]
}

// Set overflow flag
func (c *Cpu) setOf(newOf int) byte {
	if newOf != 0 {
		c.regs8[FLAG_OF] = 1
	} else {
		c.regs8[FLAG_OF] = 0
	}
	return c.regs8[FLAG_OF]
}

// Set auxiliary and overflow flag after arithmetic operations
func (c *Cpu) setAfOfArith() byte {
	c.setAf(int((c.opSource ^ c.opDest ^ uint32(c.opResult)) & 0x10))
	if c.opResult == int(c.opDest) {
		return c.setOf(0)
	}
	return c.setOf(1 & (int(c.regs8[FLAG_CF]) ^ int(c.opSource)>>(c.topBit()-1)))
}

// Assemble and return emulated CPU FLAGS register in scratch_uint
func (c *Cpu) makeFlags() {
	c.scratchUint = 0xF002
	for i := 9; i >= 0; i-- {
		c.scratchUint += uint32(c.regs8[FLAG_CF+i]) << c.biosTableLookup[TABLE_FLAGS_BITFIELDS][i]
	}
}

// Set emulated CPU FLAGS register from regs8[FLAG_xx] values
func (c *Cpu) setFlags(newFlags int) {
	for i := 9; i >= 0; i-- {
		c.regs8[FLAG_CF+i] = BoolByte((c.scratchUint >> c.biosTableLookup[TABLE_FLAGS_BITFIELDS][i]) & uint32(newFlags))
	}
}
