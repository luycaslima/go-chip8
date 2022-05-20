package emulator

import (
	"go-chip8/config"
	m "go-chip8/memory"
	r "go-chip8/registers"
)

type Chip8 struct {
	m.Memory
	r.Registers
	r.Stack
}

func (ch8 *Chip8) IsStackInBounds() {
	if ch8.SP > config.STACK_MAXSIZE {
		panic("STACK: Out of bounds ")
	}
}

func (ch8 *Chip8) PushStack(value uint16) {

	ch8.PC += 1
	ch8.IsStackInBounds()
	ch8.Stack.Stck[ch8.PC] = value

}

func (ch8 *Chip8) PopStack() uint16 {
	ch8.IsStackInBounds()
	result := ch8.Stack.Stck[ch8.PC]
	ch8.PC -= 1
	return result
}
