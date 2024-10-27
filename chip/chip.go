package chip

import (
	"chip8/utils"
	"time"
)

type Chip struct {
	pc         uint16
	index      uint16
	registers  []uint8
	memory     []uint8
	display    []bool
	stack      utils.Stack
	delayTimer uint8
	soundTimer uint8
}

// The font for each sprite
var fontValues []uint8 = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func New() *Chip {
	// Init memory
	// As by specification, the memory should only have 4KB of space
	memory := make([]uint8, 4096)

	// 0x050 - 0x09F placement for fonts seems to have been a convention, following that
	copy(memory[0x050:0x09F], fontValues)

	// Init display
	// The display used was 64 Pixels wide and 32 Tall
	display := make([]bool, 64*32)

	// Init stack
	// Could be limited to emulate old hardware better, but currently no restrictions as far stack memory.
	stack := utils.NewStack()

	// Init registers
	// CHIP-8 had 16 8-bit general use registers (V0-VF)
	registers := make([]uint8, 16)

	chip := Chip{
		pc:         0,
		index:      0,
		registers:  registers,
		memory:     memory,
		display:    display,
		stack:      stack,
		delayTimer: 0,
		soundTimer: 0,
	}

	// Set up a way for the timers to decrement each second
	go func() {
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				// Accessing the timers via chip because they may be reassigned by a program
				chip.delayTimer--
				chip.soundTimer--
			}
		}
	}()

	return &chip
}
