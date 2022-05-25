package main

import (
	"go-chip8/config"
	"go-chip8/emulator"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	chip8 := emulator.Chip8{}
	chip8.Start("roms/INVADERS")

	//SDL BASIC INITIALIZATION
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("GO! Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(config.WIDTH*config.SCREEN_MULTIPLIER),
		int32(config.HEIGHT*config.SCREEN_MULTIPLIER), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				switch t.State {
				case sdl.PRESSED:
					key := t.Keysym.Sym
					vkey := chip8.CheckMappedKeys(key)
					if vkey != 1 {
						chip8.KeyDown(int(vkey))
					}
					if key == sdl.K_ESCAPE {
						return
					}

				case sdl.RELEASED:
					key := t.Keysym.Sym
					vkey := chip8.CheckMappedKeys(key)
					if vkey != 1 {
						chip8.KeyUp(int(vkey))
					}

				}
			}
		}

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()
		renderer.SetDrawColor(255, 255, 255, 0) //seta a cor do q for desenhado na tela (Rect, Line, e Clear)

		//Desenhar os pixels na tela
		for x := 0; x < config.WIDTH; x++ {
			for y := 0; y < config.HEIGHT; y++ {
				if chip8.IsScreenSet(x, y) {
					var r sdl.Rect
					r.X = int32(x * config.SCREEN_MULTIPLIER)
					r.Y = int32(y * config.SCREEN_MULTIPLIER)
					r.W = int32(config.SCREEN_MULTIPLIER)
					r.H = int32(config.SCREEN_MULTIPLIER)
					renderer.FillRect(&r)
				}
			}
		}
		renderer.Present()

		//TODO Separate this in a gofunc?
		//timers
		if chip8.DT > 0 {
			chip8.DT = chip8.DT - 1

		}

		if chip8.ST > 0 {
			chip8.ST = chip8.ST - 1
			//Soltar um som aqui
		}

		//Execução dos comandos aqui
		chip8.GenerateOpCode()

		//TODO try use getTicks
		time.Sleep(time.Microsecond * 16700) //60hz
	}
}
