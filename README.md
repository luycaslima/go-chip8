# GO! Chip8
"GO! Chip8" is a Chip8 emulator made in Golang as a training project for practice the language, bitwise operations and the use of sdl2.

This is an adaptation of the project from the Udemy's Course of [Daniel McCarthy](https://www.udemy.com/course/creating-a-chip-8-emulator-in-c/).

---
## What is CHIP-8 ?

"CHIP-8 is an interpreted programming language, developed by Joseph Weisbecker made on his 1802 Microprocessor. It was initially used on the COSMAC VIP and Telmac 1800 8-bit microcomputers in the mid-1970s. CHIP-8 programs are run on a CHIP-8 virtual machine. It was made to allow video games to be more easily programmed for these computers. The simplicity of CHIP-8, and its long history and popularity, has ensured that CHIP-8 emulators and programs are still being made to this day." Wikipedia

---
## How to run?

Download one of the Public Domain Roms created by the community [Link](https://johnearnest.github.io/chip8Archive/)
and put inside the projects folder.

At the terminal opened at the project's folder, execute the program

```
go run main.go
```

Next, type the path to your game and then the game will execute.

---
## Controllers 

Chip 8's Keys diagram

1	2	3	C  
4	5	6	D  
7	8	9	E  
A	0	B	F  


Mapped at the keyboard

1	2	3	4  
Q	W	E	R  
A	S	D	F  
Z	X	C	V  

---

## TODO:
- [ ] Simulate Sound
- [X] Improve Framerate
- [X] Fix the Keypress Exception caused when press a key out of the mapped keys
