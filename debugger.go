package main

import (
	"fmt"
	"io"
	"os"
)

type debugger struct {
	*virtualMachine
	pos uint16
}

func NewDebugger(r io.Reader) *debugger {
	vm := NewVirtualMachine(r)

	debugger := debugger{
		virtualMachine: vm,
	}

	return &debugger
}

func (d *debugger) DebugNextInstruction() {
	nextInstruction := d.printInstruction()
	fmt.Println("Instruction:", nextInstruction)

	// d.RunNextInstruction()
}

func (d *debugger) read() uint16 {
	result := d.memory[d.pos]
	d.pos++
	return result
}

func (d *debugger) printInstruction() string {
	address := d.address
	fmt.Println("Address:", address)

	nextInstruction := d.memory[address]
	switch nextInstruction {
	case 0:
		// halt: 0
		return "halt"
	case 1:
		// set: 1 a b
		a := d.read()
		b := d.read()
		bv := d.value(b)
		return fmt.Sprintf("set <%d> %d (%d)", a, b, bv)
	case 2:
		// push: 2 a
		a := d.read()
		a = d.value(a)
		return fmt.Sprintf("push <%d>", a)
	case 3:
		// pop: 3 a
		a := d.read()
		return fmt.Sprintf("pop <%d>", a)
	case 4:
		// eq: 4 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("eq <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 5:
		// gt: 5 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("eq <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 6:
		// jmp a
		a := d.read()
		av := d.value(a)
		return fmt.Sprintf("jmp %d (%d)", a, av)
	case 7:
		// jt a b
		a := d.read()
		b := d.read()

		av := d.value(a)
		bv := d.value(b)

		return fmt.Sprintf("jt %d (%d) %d (%d)", a, av, b, bv)
	case 8:
		// jf a b
		a := d.read()
		b := d.read()

		av := d.value(a)
		bv := d.value(b)

		return fmt.Sprintf("jf %d (%d) %d (%d)", a, av, b, bv)
	case 9:
		// add: 9 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("add <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 10:
		// mult: 10 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("mult <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 11:
		// mod: 11 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("mod <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 12:
		// and: 12 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("and <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 13:
		// or: 13 a b c
		a := d.read()
		b := d.read()
		c := d.read()

		bv := d.value(b)
		cv := d.value(c)

		return fmt.Sprintf("or <%d> %d (%d) %d (%d)", a, b, bv, c, cv)
	case 14:
		// not: 14 a b
		a := d.read()
		b := d.read()

		bv := d.value(b)

		return fmt.Sprintf("not <%d> %d (%d)", a, b, bv)
	case 15:
		// rmem: 15 a b
		a := d.read()
		b := d.read()

		bv := d.value(b)
		return fmt.Sprintf("rmem <%d> %d (%d)", a, b, bv)
	case 16:
		// wmem: 16 a b
		a := d.read()
		b := d.read()

		av := d.value(a)
		bv := d.value(b)
		return fmt.Sprintf("wmem %d (%d) %d (%d)", a, av, b, bv)
	case 17:
		// call: 17 a
		a := d.read()

		// write the address of the next instruction to the stack and jump to <a>
		next := d.address
		d.push(uint16(next))

		av := d.value(a)
		return fmt.Sprintf("call %d (%d)", a, av)
	case 18:
		// ret: 18
		return "ret"
	case 19:
		// out a
		a := d.read()
		av := d.value(a)
		return fmt.Sprintf("out %d (%d) [%s]", a, av, string(av))
	case 20:
		// in: 20 a
		a := d.read()
		// read a character from the terminal and write its ascii code to <a>;
		// it can be assumed that once input starts, it will continue until a newline is encountered;
		// this means that you can safely read whole lines from the keyboard and trust that they will be fully read
		buf := make([]byte, 1)
		os.Stdin.Read(buf)
		d.setRegister(a, uint16(buf[0]))
		return fmt.Sprintf("in <%d>", a)
	case 21:
		// noop
		return "noop"
	default:
		return fmt.Sprintf("Unknown operation: %d", nextInstruction)
	}
}
