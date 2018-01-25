package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	modulo = 32768
)

type virtualMachine struct {
	memory    map[int16]int16
	registers [8]int16
	stack     []int16
	address   int
	code      []uint16
}

func (v *virtualMachine) run() {
	fmt.Println("\n--- Starting ---\n")
	for {
		instruction := v.code[v.address]
		offset := v.execute(instruction)
		if offset != 0 {
			// if offset == 0 v.address has already been set
			v.address += offset
		}
	}
	fmt.Println("\n--- Done ---\n")
}

func (v *virtualMachine) execute(instruction uint16) int {
	var offset int
	switch instruction {
	case 0:
		fmt.Println(v.address, ":", instruction)
		panic("End of the program")
	case 6:
		// jmp a
		a := v.code[v.address+1]
		v.address = int(a)
		offset = 0
	case 7:
		// jt a b
		a := v.code[v.address+1]
		b := v.code[v.address+2]

		// 32768..32775 instead mean registers 0..7
		if a >= 32768 {
			fmt.Println(a)
			panic("Register")
		}
		if a != 0 {
			// jump to b
			v.address = int(b)
			offset = 0
		} else {
			offset = 3
		}
	case 8:
		// jf a b
		a := v.code[v.address+1]
		b := v.code[v.address+2]

		if a >= 32768 {
			fmt.Println(a)
			panic("Register")
		}
		if a == 0 {
			// jump to b
			v.address = int(b)
			offset = 0
		} else {
			offset = 3
		}
	case 19:
		// out a
		a := string(v.code[v.address+1])
		fmt.Print(a)
		offset = 2
	case 21:
		// no operation
		offset = 1
	default:
		panic(fmt.Sprintf("Unknown operation: %d", instruction))
	}
	return offset
}

func main() {
	f, err := os.Open("challenge.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(b[:50])
	fmt.Println(string([]byte{87, 101, 108}))

	// Convert []byte to []uint16
	u := make([]uint16, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		v := binary.LittleEndian.Uint16(b[i : i+2])
		u[i/2] = v
	}

	fmt.Println(u[0:5])

	vm := virtualMachine{
		memory:    make(map[int16]int16),
		registers: [8]int16{0, 0, 0, 0, 0, 0, 0, 0},
		stack:     make([]int16, 0),
		address:   0,
		code:      u,
	}
	vm.run()
}
