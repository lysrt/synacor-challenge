package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	modulo = 32768
)

type virtualMachine struct {
	memory     map[uint16]int16
	registers  [8]uint16
	stack      []uint16
	address    int
	input      []uint16
	terminated bool
}

func NewVirtualMachine(r io.Reader) *virtualMachine {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	// Convert []byte to []uint16
	u := make([]uint16, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		v := binary.LittleEndian.Uint16(b[i : i+2])
		u[i/2] = v
	}

	vm := virtualMachine{
		memory:    make(map[uint16]int16),
		registers: [8]uint16{0, 0, 0, 0, 0, 0, 0, 0},
		stack:     make([]uint16, 0),
		input:     u,
	}

	return &vm
}

func (v *virtualMachine) read() uint16 {
	result := v.input[v.address]
	v.address++
	return result
}

func (v *virtualMachine) jump(address int) {
	v.address = address
}

func isRegister(a uint16) bool {
	return a >= 32768 && a <= 32775
}

func register(a uint16) uint16 {
	return a - modulo
}

func (v *virtualMachine) value(a uint16) uint16 {
	if isRegister(a) {
		return v.registers[register(a)]
	}
	return a
}

func (v *virtualMachine) setRegister(r, value uint16) {
	if !isRegister(r) {
		panic(string(r) + " is not a register")
	}
	v.registers[register(r)] = value
}

func (v *virtualMachine) RunNextInstruction() error {
	instruction := v.read()
	switch instruction {
	case 0:
		// halt: 0
		v.terminated = true
	case 1:
		// set: 1 a b
		a := v.read()
		b := v.read()
		//set register <a> to the value of <b>
		v.setRegister(a, b)
	case 2:
		// push: 2 a
		a := v.read()
		a = v.value(a)
		// push <a> onto the stack
		// v.push(val)
	case 4:
		// eq: 4 a b c
		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
		if b == c {
			v.setRegister(a, 1)
		} else {
			v.setRegister(a, 0)
		}
	case 6:
		// jmp a
		a := v.read()
		v.jump(int(a))
	case 7:
		// jt a b
		a := v.read()
		b := v.read()

		a = v.value(a)
		if a != 0 {
			// jump to b
			v.jump(int(b))
		}
	case 8:
		// jf a b
		a := v.read()
		b := v.read()

		a = v.value(a)
		if a == 0 {
			// jump to b
			v.jump(int(b))
		}
	case 9:
		// add: 9 a b c
		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// assign into <a> the sum of <b> and <c> (modulo 32768)
		sum := (b + c) % 32768
		v.setRegister(a, sum)
	case 19:
		// out a
		a := v.read()
		fmt.Print(string(a))
	case 21:
		// noop
	default:
		return fmt.Errorf("Unknown operation: %d", instruction)
	}

	return nil
}
