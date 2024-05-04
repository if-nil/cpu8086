package cpu8086

func main() {

	cpu := NewCPU()

	cpu.regs8 = cpu.mem[REGS_BASE:]
	// cpu.regs16 = cpu.mem[REGS_BASE:]
}
