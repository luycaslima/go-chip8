package registers

import (
	"go-chip8/config"
)

type Registers struct {
	V  [config.V_MAXSIZE_REGISTERS]byte
	I  uint16
	PC uint16
	SP byte
	//Para som
	ST uint8 // Sound Timer
	DT uint8 //Delay timer
}

type Stack struct {
	Stck [config.STACK_MAXSIZE]uint16
}
