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
	registers [8]uint16
	stack     []int16
	address   int
	code      []uint16
}

func (v *virtualMachine) run() {
	fmt.Print("\n--- Starting ---\n\n")
	for {
		instruction := v.code[v.address]
		offset := v.execute(instruction)
		if offset != 0 {
			// if offset == 0 v.address has already been set
			v.address += offset
		}
	}
}

func (v *virtualMachine) execute(instruction uint16) int {
	var offset int
	fmt.Printf("Executing instruction %d\n", instruction)
	switch instruction {
	case 0:
		fmt.Println(v.address, ":", instruction)
		panic("End of the program")
	case 1:
		// set: 1 a b
		a := v.code[v.address+1]
		b := v.code[v.address+2]

		if a < 32768 || a > 32775 {
			panic(string(a) + " is not a register")
		}

		//set register <a> to the value of <b>
		registerNumber := a - 32768
		v.registers[registerNumber] = b

		offset = 3
	case 2:
		// push: 2 a
		a := v.code[v.address+1]

		if a >= 32768 {
			registerNumber := a - 32768
			registerValue := v.registers[registerNumber]
			a = registerValue
		}
		// push <a> onto the stack
		v.push(a)

		offset = 2
	case 4:
		// eq: 4 a b c
		a := v.code[v.address+1]
		b := v.code[v.address+2]
		c := v.code[v.address+3]

		if a < 32768 || a > 32775 {
			panic(string(a) + " is not a register")
		}

		if b >= 32768 {
			registerNumber := b - 32768
			registerValue := v.registers[registerNumber]
			b = registerValue
		}

		if c >= 32768 {
			registerNumber := c - 32768
			registerValue := v.registers[registerNumber]
			c = registerValue
		}

		// set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
		if b == c {
			v.registers[a-32768] = 1
		} else {
			v.registers[a-32768] = 0
		}

		offset = 4
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
			registerNumber := a - 32768
			registerValue := v.registers[registerNumber]
			a = registerValue
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
			registerNumber := a - 32768
			registerValue := v.registers[registerNumber]
			a = registerValue
		}
		if a == 0 {
			// jump to b
			v.address = int(b)
			offset = 0
		} else {
			offset = 3
		}
	case 9:
		// add: 9 a b c
		a := v.code[v.address+1]
		b := v.code[v.address+2]
		c := v.code[v.address+3]

		if a < 32768 || a > 32775 {
			panic(string(a) + " is not a register")
		}

		if b >= 32768 {
			registerNumber := b - 32768
			registerValue := v.registers[registerNumber]
			b = registerValue
		}

		if c >= 32768 {
			registerNumber := c - 32768
			registerValue := v.registers[registerNumber]
			c = registerValue
		}

		sum := (b + c) % 32768

		v.registers[a-32768] = sum

		offset = 4
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
		registers: [8]uint16{0, 0, 0, 0, 0, 0, 0, 0},
		stack:     make([]int16, 0),
		address:   0,
		code:      u,
	}
	vm.run()
}
