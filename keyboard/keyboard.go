package keyboard

import "github.com/veandco/go-sdl2/sdl"

/*Chip 8's Keys diagram
1	2	3	C
4	5	6	D
7	8	9	E
A	0	B	F
*/

/*Mapped keyboard
1	2	3	4
Q	W	E	R
A	S	D	F
Z	X	C	V
*/

var MappedKeys = [16]sdl.Keycode{
	sdl.K_1, sdl.K_2, sdl.K_3, sdl.K_4,
	sdl.K_q, sdl.K_w, sdl.K_e, sdl.K_r,
	sdl.K_a, sdl.K_s, sdl.K_d, sdl.K_f,
	sdl.K_z, sdl.K_x, sdl.K_c, sdl.K_v}

type Keyboard struct {
	Keys [16]bool
}

//Make the virtual key released false on the MappedKeys
func (keyboard *Keyboard) KeyUp(key int) {
	keyboard.Keys[key] = false
}

//Make the virtual key released true on the MappedKeys
func (keyboard *Keyboard) KeyDown(key int) {
	keyboard.Keys[key] = true
}

//Return if the key is still pressed or not
func (keyboard *Keyboard) KeyIsPressed(keyIndex int) bool {
	return keyboard.Keys[keyIndex]
}

func (keyboard *Keyboard) CheckMappedKeys(key sdl.Keycode) int {
	for i := 0; i < 16; i++ {
		if key == MappedKeys[i] {
			return i
		}
	}
	return -1
}

func (keyboard *Keyboard) WaitKeyEvent() uint8 {
	for {
		ev := sdl.WaitEvent()
		if ev.GetType() != sdl.KEYDOWN {
			continue
		}
		switch t := ev.(type) {
		case *sdl.KeyboardEvent:
			key := t.Keysym.Sym
			vkey := keyboard.CheckMappedKeys(key)
			if vkey != -1 {
				return uint8(vkey)
			}
		}
	}
}
