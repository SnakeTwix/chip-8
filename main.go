package main

import (
	"chip8/chip"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	ibm, err := os.ReadFile("chip_binaries/2-ibm-logo.ch8")
	if err != nil {
		log.Fatal(err)
	}

	workingChip := chip.NewWithMemory(ibm)
	outputChannel := workingChip.GetOutputChannel()

	go func() {
		for {
			fmt.Println(<-outputChannel)
		}
	}()

	for {
		workingChip.Run()
		time.Sleep(time.Millisecond * 100)
	}

}
