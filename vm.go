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
	memory     map[uint16]uint16
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
		memory:    make(map[uint16]uint16),
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

func (v *virtualMachine) push(value uint16) {
	v.stack = append(v.stack, value)
}

func (v *virtualMachine) pop() uint16 {
	if len(v.stack) == 0 {
		panic("empty stack")
	}
	result := v.stack[len(v.stack)-1]
	v.stack = v.stack[:len(v.stack)-1]
	return result
}

func (v *virtualMachine) readMemory(address uint16) uint16 {
	return v.memory[address]
}

func (v *virtualMachine) writeMemory(address, value uint16) {
	v.memory[address] = value
}

func (v *virtualMachine) RunNextInstruction() error {
	instruction := v.read()
	// fmt.Println("Executing instruction", instruction)
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
		v.push(a)
	case 3:
		// pop: 3 a
		a := v.read()
		// remove the top element from the stack and write it into <a>; empty stack = error
		value := v.pop()
		v.setRegister(a, value)
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
	case 5:
		// gt: 5 a b c
		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
		if b > c {
			v.setRegister(a, 1)
		} else {
			v.setRegister(a, 0)
		}
	case 6:
		// jmp a
		a := v.read()
		// jump to <a>
		a = v.value(a)
		v.jump(int(a))
	case 7:
		// jt a b
		a := v.read()
		b := v.read()

		a = v.value(a)
		if a != 0 {
			// if <a> is nonzero, jump to <b>
			b = v.value(b)
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
	case 10:
		// mult: 10 a b c
		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// store into <a> the product of <b> and <c> (modulo 32768)
		product := (b * c) % 32768
		v.setRegister(a, product)
	case 11:
		// mod: 11 a b c
		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// store into <a> the remainder of <b> divided by <c>
		remainder := b % c
		v.setRegister(a, remainder)
	case 12:
		// and: 12 a b c
		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// stores into <a> the bitwise and of <b> and <c>
		sum := b & c
		v.setRegister(a, sum)
	case 13:
		// or: 13 a b c

		a := v.read()
		b := v.read()
		c := v.read()

		b = v.value(b)
		c = v.value(c)

		// stores into <a> the bitwise or of <b> and <c>
		sum := b | c
		v.setRegister(a, sum)
	case 14:
		// not: 14 a b
		a := v.read()
		b := v.read()

		b = v.value(b)
		// stores 15-bit bitwise inverse of <b> in <a>
		not := ^b & 0x7fff
		v.setRegister(a, not)
	case 15:
		// rmem: 15 a b
		a := v.read()
		b := v.read()

		// read memory at address <b> and write it to <a>
		b = v.value(b)
		value := v.readMemory(b)
		v.setRegister(a, value)
	case 17:
		// call: 17 a
		a := v.read()

		// write the address of the next instruction to the stack and jump to <a>
		next := v.address
		v.push(uint16(next))

		a = v.value(a)
		v.jump(int(a))
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
