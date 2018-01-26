package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("challenge.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	vm := NewVirtualMachine(f)

	log.Println("Starting VM...")
	fmt.Println()

	for !vm.terminated {
		err := vm.RunNextInstruction()
		if err != nil {
			log.Fatalf("fatal vm error: %q", err)
			break
		}
	}

	log.Println("VM terminated successfully")
}
