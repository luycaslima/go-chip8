package main

import (
	"bufio"
	"fmt"
	"go-chip8/config"
	"go-chip8/emulator"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 160, 144 //Native 160x144 640, 576

func loadGame(path string) ([]byte, error) {
	infile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	stats, statsErr := infile.Stat() //retornada a estrutura e detalhes de um arquivo se houver algum problema retorna patherror
	if statsErr != nil {
		return nil, statsErr
	}

	size := stats.Size() //tamanho em bytes do arquivo
	bytes := make([]byte, size)

	bufr := bufio.NewReader(infile)
	_, err = bufr.Read(bytes)
	return bytes, err
}

func main() {
	chip8 := emulator.Chip8{}
	chip8.Registers.SP = 0

	chip8.PushStack(0xff)
	chip8.PushStack(0xaa)

	fmt.Printf("%x\n", chip8.PopStack())
	fmt.Printf("%x\n", chip8.PopStack())

	//fmt.Println("byteArray: ", rom)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit() //defer chama essa função ao fim do main

	//WINDOW
	window, err := sdl.CreateWindow("GO! Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(config.WINDOW_WIDTH*config.WINDOW_MULTIPLIER), int32(config.WINDOW_HEIGHT*config.WINDOW_MULTIPLIER), sdl.WINDOW_SHOWN) //retorna multiplas coisas
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	//RENDER
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err) //Encerra o programa e retorna o erro gerado
	}
	defer renderer.Destroy()

	//Texture
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight)) //TODO: Ver oq sao essas flags
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4) //Mutltiplicado por 4 pois cada pixel armaze 4 bytes(um para alpha e os outros pro RGB)

	var frameStart time.Time
	var elapsedTime float32
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				return
			}
		}
		tex.Update(nil, pixels, winWidth*4) //pitch é basicamente quantos bytes tem o width da janela
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime*1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}

}
