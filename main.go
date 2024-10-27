package main

import (
	"chip8/chip"
	"fmt"
	"time"
)

func main() {
	workingChip := chip.New()

	time.Sleep(time.Second * 10)

	fmt.Println(workingChip)
}
