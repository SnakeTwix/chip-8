package chip

import (
	"chip8/utils"
	"math"
	"strings"
	"time"
)

type Chip struct {
	pc            uint16
	index         uint16
	registers     []uint8
	memory        []uint8
	display       []bool
	stack         utils.Stack
	delayTimer    uint8
	soundTimer    uint8
	outputChannel chan string
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

	// The output channel
	// Not ever specified anywhere, just for me to be able to read data and output it to wherever I want
	output := make(chan string)

	chip := Chip{
		pc:            512,
		index:         0,
		registers:     registers,
		memory:        memory,
		display:       display,
		stack:         stack,
		delayTimer:    0,
		soundTimer:    0,
		outputChannel: output,
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

func NewWithMemory(memory []uint8) *Chip {
	chip := New()

	copy(chip.memory[512:], memory)

	return chip
}

func (c *Chip) Run() {
	instruction := c.fetch()
	c.decode(instruction)

	//fmt.Println("Memory: ", c.memory)
	//fmt.Println("Display: ", c.display)
}

func (c *Chip) fetch() uint16 {
	byteOne := c.memory[c.pc]
	byteTwo := c.memory[c.pc+1]
	c.pc += 2

	var instruction = uint16(byteOne) << 8
	instruction += uint16(byteTwo)

	return instruction
}

func (c *Chip) decode(instruction uint16) {
	switch instruction & 0xF000 {
	case 0x0000:
		if instruction&0x0FFF == 0x00E0 {
			for index := range c.display {
				c.display[index] = false
			}
		} else if instruction == 0x00EE {
			c.pc = c.stack.Pop()
		}
	case 0x1000:
		c.pc = instruction & 0x0FFF
	case 0x2000:
		c.stack.Push(c.pc)
		c.pc = instruction & 0x0FFF
	case 0x3000:
		registerAddress := (instruction & 0x0F00) >> 8
		value := uint8(instruction & 0x00FF)

		if value == c.registers[registerAddress] {
			c.pc += 2
		}
	case 0x4000:
		registerAddress := (instruction & 0x0F00) >> 8
		value := uint8(instruction & 0x00FF)

		if value != c.registers[registerAddress] {
			c.pc += 2
		}
	case 0x5000:
		registerXAddress := (instruction & 0x0F00) >> 8
		registerYAddress := (instruction & 0x00F0) >> 4

		if c.registers[registerXAddress] == c.registers[registerYAddress] {
			c.pc += 2
		}
	case 0x6000:
		registerAddress := (instruction & 0x0F00) >> 8
		value := uint8(instruction & 0x00FF)

		c.registers[registerAddress] = value
	case 0x7000:
		registerAddress := (instruction & 0x0F00) >> 8
		value := uint8(instruction & 0x00FF)

		c.registers[registerAddress] += value
	case 0x8000:
		registerXAddress := (instruction & 0x0F00) >> 2
		registerYAddress := (instruction & 0x00F0) >> 1

		switch instruction & 0x000F {
		case 0x0000:
			c.registers[registerXAddress] = c.registers[registerYAddress]
		case 0x0001:
			c.registers[registerXAddress] = c.registers[registerXAddress] | c.registers[registerYAddress]
		case 0x0002:
			c.registers[registerXAddress] = c.registers[registerXAddress] & c.registers[registerYAddress]
		case 0x0003:
			c.registers[registerXAddress] = c.registers[registerXAddress] ^ c.registers[registerYAddress]
		case 0x0004:
			if c.registers[registerXAddress] > math.MaxUint8-c.registers[registerYAddress] {
				c.registers[0x000F] = 1
			} else {
				c.registers[0x000F] = 0
			}

			c.registers[registerXAddress] += c.registers[registerYAddress]
		case 0x0007:
			if c.registers[registerYAddress] > c.registers[registerXAddress] {
				c.registers[0x000F] = 1
			} else {
				c.registers[0x000F] = 0
			}

			c.registers[registerXAddress] = c.registers[registerYAddress] - c.registers[registerXAddress]
		}
	case 0x9000:
		registerXAddress := (instruction & 0x0F00) >> 2
		registerYAddress := uint8(instruction&0x00F0) >> 1

		if c.registers[registerXAddress] != c.registers[registerYAddress] {
			c.pc += 2
		}
	case 0xA000:
		c.index = instruction & 0x0FFF
	case 0xD000:
		registerXAddress := (instruction & 0x0F00) >> 8
		registerYAddress := (instruction & 0x00F0) >> 4
		yCord := uint16(c.registers[registerYAddress]) % 32

		c.registers[0x000F] = 0

		amountOfRows := instruction & 0x000F
		for i := range amountOfRows {
			if yCord >= 32 {
				break
			}
			xCord := uint16(c.registers[registerXAddress]) % 64
			spriteByte := c.memory[c.index+i]

			for bit := 7; bit >= 0; bit-- {
				if xCord >= 64 {
					break
				}
				currentBit := (spriteByte & (1 << bit)) >> bit

				if c.display[yCord*64+xCord] == true && currentBit == 0 {
					c.registers[0x000F] = 1
				}

				c.display[yCord*64+xCord] = currentBit == 1
				xCord++
			}

			yCord++
		}

		c.RenderDisplay()
	}

}

func (c *Chip) GetOutputChannel() <-chan string {
	return c.outputChannel
}

func (c *Chip) RenderDisplay() {
	var builder strings.Builder

	for i := range c.display {
		if i != 0 && i%64 == 0 {
			builder.WriteString("\n")
		}

		if c.display[i] == true {
			builder.WriteString("██")
		} else {
			builder.WriteString("  ")
		}

	}

	c.outputChannel <- builder.String()
}
