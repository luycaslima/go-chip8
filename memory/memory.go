package memory

import (
	ct "go-chip8/config"
)

type Memory struct {
	Memory [ct.MEMORY_SIZE]uint8
}

//TODO: Adicionar um retorno de possivel error?
func (mem *Memory) IsInMemoryBounds(index int) {
	if index <= 0 || index > ct.MEMORY_SIZE {
		panic("Memory out of bounds")
	}
}

func (mem *Memory) MemorySet(index int, value uint8) {

	mem.IsInMemoryBounds(index)
	mem.Memory[index] = value
}

func (mem *Memory) MemoryGet(index int) uint8 {
	mem.IsInMemoryBounds(index)
	return mem.Memory[index]
}
