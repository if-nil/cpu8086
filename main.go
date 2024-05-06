package cpu8086

import (
	"log"
	"os"
)

func main() {

	cpu := NewCPU()

	cpu.regs16.Set(REG_CS, 0xF000)

	// Trap flag off
	cpu.regs8.Set(FLAG_TF, 0)

	// Set DL equal to the boot device: 0 for the FD, or 0x80 for the HD. Normally, boot from the FD.
	//  But, if the HD image file is prefixed with @, then boot from the HD
	cpu.regs8.Set(REG_DL, 0)

	// Open BIOS (file id disk[2]), floppy disk image (disk[1]), and hard disk image (disk[0]) if specified
	var err error
	for i := 3; i > 0; i-- {
		cpu.disk[i-1], err = os.OpenFile(os.Args[i], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error opening file: %v, err: %v", os.Args[i], err)
		}
	}

	// Set CX:AX equal to the hard disk image size, if present
	if cpu.disk[0] != nil {
		ret, err := cpu.disk[0].Seek(0, 2)
		if err != nil {
			log.Fatalf("Error seeking: %v", err)
		}
		cpu.regs16.Set(REG_AX, uint16(ret>>9>>16))
		cpu.regs16.Set(REG_CX, uint16(ret>>9&0xFFFF))
	}
}
