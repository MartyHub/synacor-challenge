package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Number uint16

const modulo = 32768

type VirtualMachine struct {
	memory      [32767]Number
	registers   [8]Number
	stack       []Number
	memoryIndex int
	input       []rune
}

func NewVirtualMachine(program []Number) *VirtualMachine {
	result := &VirtualMachine{}

	copy(result.memory[:], program)

	return result
}

func (vm *VirtualMachine) Run() {
	for {
		opcode := vm.toValue(vm.current())

		switch opcode {
		case 0:
			vm.halt()
		case 1:
			vm.set()
		case 2:
			vm.push()
		case 3:
			vm.pop()
		case 4:
			vm.eq()
		case 5:
			vm.gt()
		case 6:
			vm.jump()
		case 7:
			vm.jt()
		case 8:
			vm.jf()
		case 9:
			vm.add()
		case 10:
			vm.mul()
		case 11:
			vm.mod()
		case 12:
			vm.and()
		case 13:
			vm.or()
		case 14:
			vm.not()
		case 15:
			vm.read()
		case 16:
			vm.write()
		case 17:
			vm.call()
		case 18:
			vm.ret()
		case 19:
			vm.print()
		case 20:
			vm.ask()
		case 21:
			vm.noop()
		default:
			log.Fatalf("Unknown opcode %d from value %d at address %d", opcode, vm.current(), vm.memoryIndex)
		}

		vm.next()
	}
}

func (vm *VirtualMachine) current() Number {
	return vm.memory[vm.memoryIndex]
}

func (vm *VirtualMachine) next() Number {
	vm.memoryIndex++

	return vm.memory[vm.memoryIndex]
}

func (vm *VirtualMachine) toRegisterIndex(n Number) int {
	return int(n) - modulo
}

func (vm *VirtualMachine) toValue(n Number) Number {
	if n >= 32776 {
		log.Fatalf("Invalid value %d", n)
	}

	if n >= modulo {
		return vm.registers[vm.toRegisterIndex(n)]
	}

	return n
}

func (vm *VirtualMachine) and() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())
	vm.registers[registerIndex] = (value1 & value2) % modulo
}

func (vm *VirtualMachine) or() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())
	vm.registers[registerIndex] = (value1 | value2) % modulo
}

func (vm *VirtualMachine) not() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value := vm.toValue(vm.next())
	vm.registers[registerIndex] = (^value) % modulo
}

func (vm *VirtualMachine) add() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())
	vm.registers[registerIndex] = (value1 + value2) % modulo
}

func (vm *VirtualMachine) mul() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())
	vm.registers[registerIndex] = (value1 * value2) % modulo
}

func (vm *VirtualMachine) mod() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())
	vm.registers[registerIndex] = (value1 % value2) % modulo
}

func (vm *VirtualMachine) gt() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())

	if value1 > value2 {
		vm.registers[registerIndex] = 1
	} else {
		vm.registers[registerIndex] = 0
	}
}

func (vm *VirtualMachine) eq() {
	registerIndex := vm.toRegisterIndex(vm.next())
	value1 := vm.toValue(vm.next())
	value2 := vm.toValue(vm.next())

	if value1 == value2 {
		vm.registers[registerIndex] = 1
	} else {
		vm.registers[registerIndex] = 0
	}
}

func (vm *VirtualMachine) halt() {
	os.Exit(0)
}

func (vm *VirtualMachine) jump() {
	vm.memoryIndex = int(vm.toValue(vm.next())) - 1
}

func (vm *VirtualMachine) jf() {
	if vm.toValue(vm.next()) == 0 {
		vm.jump()
	} else {
		vm.next()
	}
}

func (vm *VirtualMachine) jt() {
	if vm.toValue(vm.next()) != 0 {
		vm.jump()
	} else {
		vm.next()
	}
}

func (vm *VirtualMachine) noop() {
}

func (vm *VirtualMachine) print() {
	fmt.Print(string(rune(vm.toValue(vm.next()))))
}

func (vm *VirtualMachine) ask() {
	if len(vm.input) == 0 {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		vm.input = []rune(line)
	}

	registerIndex := vm.toRegisterIndex(vm.next())

	vm.registers[registerIndex] = Number(vm.input[0])
	vm.input = vm.input[1:]
}

func (vm *VirtualMachine) pop() {
	l := len(vm.stack)

	if l == 0 {
		panic("Stack underflow")
	}

	value := vm.stack[l-1]
	vm.stack = vm.stack[:l-1]

	registerIndex := vm.toRegisterIndex(vm.next())
	vm.registers[registerIndex] = value
}

func (vm *VirtualMachine) push() {
	vm.stack = append(vm.stack, vm.toValue(vm.next()))
}

func (vm *VirtualMachine) read() {
	registerIndex := vm.toRegisterIndex(vm.next())
	vm.registers[registerIndex] = vm.memory[vm.toValue(vm.next())]
}

func (vm *VirtualMachine) write() {
	memoryIndex := vm.toValue(vm.next())
	vm.memory[memoryIndex] = vm.toValue(vm.next())
}

func (vm *VirtualMachine) call() {
	vm.stack = append(vm.stack, Number(vm.memoryIndex+2))
	vm.jump()
}

func (vm *VirtualMachine) ret() {
	l := len(vm.stack)

	if l == 0 {
		vm.halt()
	} else {
		vm.memoryIndex = int(vm.stack[l-1]) - 1
		vm.stack = vm.stack[:l-1]
	}
}

func (vm *VirtualMachine) set() {
	registerIndex := vm.toRegisterIndex(vm.next())
	vm.registers[registerIndex] = vm.toValue(vm.next())
}
