package main

import (
	"encoding/binary"
	"log"
	"os"
)

func readFile(fileName string) []byte {
	result, err := os.ReadFile(fileName)

	if err != nil {
		log.Fatalf("Failed to read %v: %v", fileName, err)
	}

	return result
}

func readProgram(b []byte) []Number {
	var result []Number

	for i := 0; i < len(b); i += 2 {
		result = append(result, Number(binary.LittleEndian.Uint16(b[i:i+2])))
	}

	return result
}

func main() {
	NewVirtualMachine(readProgram(readFile("challenge.bin"))).Run()
}
