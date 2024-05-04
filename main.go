package cpu8086

func main() {

	cpu := NewCPU()

	cpu.regs16.Set(REG_CS, 0xF000)

	// Trap flag off
	cpu.regs8.Set(FLAG_TF, 0)
}
