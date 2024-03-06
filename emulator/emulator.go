package emulator

/*
Technical Reference  CHIP 8
http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
*/

import (
	"fmt"
	"go-chip8/config"
	k "go-chip8/keyboard"
	"math/rand"
	"os"
)

// FONT Characters  in hexadecimal
var DefaultCharacterSet = [...]byte{
	0xf0, 0x90, 0x90, 0x90, 0xf0,
	0x20, 0x60, 0x20, 0x20, 0x70,
	0xf0, 0x10, 0xf0, 0x80, 0xf0,
	0xf0, 0x10, 0xf0, 0x10, 0xf0,
	0x90, 0x90, 0xf0, 0x10, 0x10,
	0xf0, 0x80, 0xf0, 0x10, 0xf0,
	0xf0, 0x80, 0xf0, 0x90, 0xf0,
	0xf0, 0x10, 0x20, 0x40, 0x40,
	0xf0, 0x90, 0xf0, 0x90, 0xf0,
	0xf0, 0x90, 0xf0, 0x10, 0xf0,
	0xf0, 0x90, 0xf0, 0x90, 0x90,
	0xe0, 0x90, 0xe0, 0x90, 0xe0,
	0xf0, 0x80, 0x80, 0x80, 0xf0,
	0xe0, 0x90, 0x90, 0x90, 0xe0,
	0xf0, 0x80, 0xf0, 0x80, 0xf0,
	0xf0, 0x80, 0xf0, 0x80, 0x80,
} //the "..." garantee that the array will be created with the fixed sized of the content

type Chip8 struct {
	Memory     [config.MEMORY_SIZE]byte //Memory
	k.Keyboard                          //Keyboard
	Screen     [config.WIDTH][config.HEIGHT]bool
	V          [16]uint8  //V Registers
	Stack      [16]uint16 //Stack
	DT         uint8      //Delay Timer
	ST         uint8      //Sound Timer
	SP         uint8      //Stack Pointer

	I  uint16 //FLAG
	PC uint16 //PROGRAM COUNTER

	//Instructions
	opcode uint16 //operation instruction
	nnn    uint16 //nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
	n      uint8  //n or nibble - A 4-bit value, the lowest 4 bits of the instruction
	x      uint8  //x - A 4-bit value, the lower 4 bits of the high byte of the instruction
	y      uint8  //y - A 4-bit value, the upper 4 bits of the low byte of the instruction
	kk     uint8  //kk or byte - An 8-bit value, the lowest 8 bits of the instruction
	Speed  int
}

func (ch8 *Chip8) Start(romPath string) {
	//Store the Character set at the 0x00
	copy(ch8.Memory[config.CHARACTER_SET_ADDRESS:], DefaultCharacterSet[:])
	//Load the programm
	ch8.loadRom(romPath)

	//fmt.Printf("%x \n", ch8.Memory)

	//Set the start the program at the PC
	ch8.PC = config.PROGRAM_LOAD_ADDRESS
}

func (ch8 *Chip8) loadRom(path string) {
	rom, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	copy(ch8.Memory[config.PROGRAM_LOAD_ADDRESS:], rom) //Load the array of bytes at the 0x200 position of the memory
}

// Check if the pixel of the Screen is lit or not
func (ch8 *Chip8) IsScreenSet(x, y int) bool {
	return ch8.Screen[x][y]
}

func (ch8 *Chip8) pushStack() {
	ch8.SP = ch8.SP + 1
	if ch8.SP > 16 {
		panic("Stack out of bounds")
	}
	ch8.Stack[ch8.SP] = ch8.PC
	ch8.PC = ch8.nnn

}

func (ch8 *Chip8) popStack() uint16 {
	if ch8.SP > 16 {
		panic("Stack out of bounds")
	}
	result := ch8.Stack[ch8.SP]
	ch8.SP = ch8.SP - 1
	return result
}

func (ch8 *Chip8) GenerateOpCode() {

	highByte := ch8.Memory[int(ch8.PC)]
	lowerByte := ch8.Memory[int(ch8.PC)+1]

	ch8.opcode = (uint16(highByte) << 8) | uint16(lowerByte)
	//fmt.Printf("opcode: %x \n", ch8.opcode)
	ch8.nnn = ch8.opcode & 0x0FFF
	ch8.n = lowerByte & 0x0F
	ch8.x = highByte & 0x0F
	ch8.y = (lowerByte >> 4) & 0x0F
	ch8.kk = lowerByte

	ch8.PC += 2
	err := ch8.executeOpCode()

	if err != nil {
		panic(err)
	}

}

func (ch8 *Chip8) executeOpCode() error {

	switch ch8.opcode & 0xF000 {
	case 0x0000:
		switch ch8.kk {
		case 0xE0: //CLS - Clear Screen
			ch8.Screen = [64][32]bool{}
		case 0xEE: // RET - Return from the subroutine
			ch8.PC = ch8.popStack()
		default:
			return fmt.Errorf("unknown opcode %x", ch8.opcode)
		}
	case 0x1000: //JP - jump to  the nnn location
		ch8.PC = ch8.nnn
	case 0x2000: // CALL addr Call subroutine at nnn.
		ch8.pushStack()
	case 0x3000: //3xkk - SE Vx, byte Skip next instruction if Vx = kk.
		if ch8.V[ch8.x] == ch8.kk {
			ch8.PC = ch8.PC + 2
		}
	case 0x4000: //4xkk - SNE Vx, byte Skip next instruction if Vx != kk.
		if ch8.V[ch8.x] != ch8.kk {
			ch8.PC = ch8.PC + 2
		}
	case 0x5000: //5xy0 - SE Vx, Vy Skip next instruction if Vx = Vy.
		if ch8.V[ch8.x] == ch8.V[ch8.y] {
			ch8.PC = ch8.PC + 2
		}
	case 0x6000: //6xkk - LD Vx, byte Set Vx = kk.
		ch8.V[ch8.x] = ch8.kk
	case 0x7000: //7xkk - ADD Vx, byte Set Vx = Vx + kk.
		ch8.V[ch8.x] = ch8.V[ch8.x] + ch8.kk
	case 0x8000:
		switch ch8.n {
		case 0x00: // LD Vx, Vy Set Vx = Vy.
			ch8.V[ch8.x] = ch8.V[ch8.y]
		case 0x01: //8xy1 - OR Vx, Vy Set Vx = Vx OR Vy.
			ch8.V[ch8.x] = ch8.V[ch8.x] | ch8.V[ch8.y]
		case 0x02: //8xy2 - AND Vx, Vy Set Vx = Vx AND Vy.
			ch8.V[ch8.x] = ch8.V[ch8.x] & ch8.V[ch8.y]
		case 0x03: //8xy3 - XOR Vx, Vy Set Vx = Vx XOR Vy.
			ch8.V[ch8.x] = ch8.V[ch8.x] ^ ch8.V[ch8.y]
		case 0x04: //8xy4 - ADD Vx, Vy Set Vx = Vx + Vy, set VF = carry.
			sum := uint16(ch8.V[ch8.x] + ch8.V[ch8.y])
			if sum > 255 {
				ch8.V[0x0F] = 1 //Carry
			} else {
				ch8.V[0x0F] = 0 //Carry
			}
			ch8.V[ch8.x] = byte(sum) //This ignores the high bytes and only storess the 8 lowest?
		case 0x05: //8xy5 - SUB Vx, Vy Set Vx = Vx - Vy, set VF = NOT borrow.
			ch8.V[0xF] = 0 //Carry
			if ch8.V[ch8.x] > ch8.V[ch8.y] {
				ch8.V[0xF] = 1 //Carry
			}
			ch8.V[ch8.x] = ch8.V[ch8.x] - ch8.V[ch8.y]
		case 0x06: //8xy6 - SHR Vx {, Vy} Set Vx = Vx SHR 1.
			ch8.V[0xF] = ch8.V[ch8.x] & 0x01
			ch8.V[ch8.x] = ch8.V[ch8.x] / 2
		case 0x07: //8xy7 - SUBN Vx, Vy Set Vx = Vy - Vx, set VF = NOT borrow.
			ch8.V[0xF] = 0 //Carry
			if ch8.V[ch8.y] > ch8.V[ch8.x] {
				ch8.V[0xF] = 1 //Carry
			}
			ch8.V[ch8.x] = ch8.V[ch8.y] - ch8.V[ch8.x]
		case 0x0E: //8xyE - SHL Vx {, Vy} Set Vx = Vx SHL 1.
			//If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
			ch8.V[0x0f] = ch8.V[ch8.x] & 0b10000000
			ch8.V[ch8.x] = ch8.V[ch8.x] * 2
		default:
			return fmt.Errorf("unknown opcode %x", ch8.opcode)
		}
	case 0x9000: //9xy0 - SNE Vx, Vy Skip next instruction if Vx != Vy.
		if ch8.V[ch8.x] != ch8.V[ch8.y] {
			ch8.PC = ch8.PC + 2
		}
	case 0xA000: //LD I, addr Set I = nnn.
		ch8.I = ch8.nnn
	case 0xB000: //JP V0, addr Jump to location nnn + V0.
		ch8.PC = uint16(ch8.V[0x0]) + ch8.nnn
	case 0xC000: //Cxkk - RND Vx, byte Set Vx = random byte AND kk.
		ch8.V[ch8.x] = uint8(rand.Intn(256)) & ch8.kk
	case 0xD000: //DRW Vx, Vy, nibble Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		/*	The interpreter reads n bytes from memory, starting at the address stored in I.
			These bytes are then displayed as sprites on Screen at coordinates (Vx, Vy). Sprites are XORed onto the existing Screen.
			If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.
			If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the Screen.
			See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more information on the Chip-8 Screen and sprites.*/
		x := int(ch8.V[ch8.x])
		y := int(ch8.V[ch8.y])
		ch8.V[0xF] = 0
		for ly := 0; ly < int(ch8.n); ly++ { // Big O(8n)
			sprByte := ch8.Memory[ch8.I+uint16(ly)]
			for lx := 0; lx < 8; lx++ { //percorre os 8 bits
				if (sprByte & (0b10000000 >> lx)) == 0 { //If this bit is 0, ignore it
					continue
				}
				isPixelLit := ch8.Screen[(lx+x)%config.WIDTH][(ly+y)%config.HEIGHT]
				if isPixelLit { //If the pixel is already lit, set the flag to 1
					ch8.V[0xF] = 1
				}
				ch8.Screen[(lx+int(x))%config.WIDTH][(ly+int(y))%config.HEIGHT] = (isPixelLit || true) && !(isPixelLit && true) //XORed the value
			}
		}
	case 0xE000:
		switch ch8.kk {
		// Ex9e - SKP Vx, Skip the next instruction if the key with the value of Vx is pressed
		case 0x9e:
			if ch8.KeyIsPressed(int(ch8.V[ch8.x])) {
				ch8.PC = ch8.PC + 2
			}
		// Exa1 - SKNP Vx - Skip the next instruction if the key with the value of Vx is not pressed
		case 0xa1:
			if !ch8.KeyIsPressed(int(ch8.V[ch8.x])) {
				ch8.PC = ch8.PC + 2
			}
		default:
			return fmt.Errorf("unknown opcode %x", ch8.opcode)
		}
	case 0xF000:
		switch ch8.kk {
		case 0x07: //LD Vx, DT Set Vx = delay timer value.
			ch8.V[ch8.x] = ch8.DT
		case 0x0A: // LD Vx, K Wait for a key press, store the value of the key in Vx.
			pressedKey := ch8.WaitKeyEvent()
			ch8.V[ch8.x] = pressedKey
		case 0x15: //Fx15 - LD DT, Vx Set delay timer = Vx.
			ch8.DT = ch8.V[ch8.x]
		case 0x18: //LD ST, Vx Set sound timer = Vx.
			ch8.ST = ch8.V[ch8.x]
		case 0x1E: //Fx1E - ADD I, Vx Set I = I + Vx.
			ch8.I = ch8.I + uint16(ch8.V[ch8.x])
		case 0x29: // Fx29- LD F, Vx Set I = location of sprite for digit Vx.
			ch8.I = uint16(ch8.V[ch8.x]) * 5
		case 0x33: //Fx33 - LD B, Vx Store BCD representation of Vx in memory locations I, I+1, and I+2.
			//The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I
			//, the tens digit at location I+1, and the ones digit at location I+2.
			ch8.Memory[ch8.I] = uint8((uint16(ch8.V[ch8.x]) % 1000) / 100)
			ch8.Memory[ch8.I+1] = (ch8.V[ch8.x] / 10) % 10
			ch8.Memory[ch8.I+2] = ch8.V[ch8.x] % 10
		case 0x55: //LD [I], Vx Store registers V0 through Vx in memory starting at location I.
			for i := 0; i <= int(ch8.x); i++ { //V total registers
				ch8.Memory[ch8.I+uint16(i)] = ch8.V[i]
			}
		case 0x65: // LD Vx, [I] Read registers V0 through Vx from memory starting at location I.
			for i := 0; i <= int(ch8.x); i++ { //V total registers
				ch8.V[i] = ch8.Memory[ch8.I+uint16(i)]
			}
		default:
			return fmt.Errorf("unknown opcode %x", ch8.opcode)
		}
	default:
		return fmt.Errorf("unknown opcode %x", ch8.opcode)
	}

	return nil
}
