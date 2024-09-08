package main

import (
	"ash/compiler"
	"ash/lexer"
	"ash/parser"
	"ash/repl"
	"ash/utils"
	"ash/vm"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args

	switch len(args) {
	case 1:
		fmt.Println("Ash 0.0.1")
		repl.Start(os.Stdin, os.Stdout)

	case 2:
		if args[1] == "-v" || args[1] == "--version" {
			fmt.Println("Ash 0.0.1")
		} else {
			run(args[1])
		}

	}
}

func run(filename string) {
	if !strings.HasSuffix(filename, ".ash") {
		filename += ".ash"
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Errorf("could not read %s: %v", filename, err)
		os.Exit(1)
	}

	p := parser.New(lexer.New(string(data)))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		utils.PrintParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	c := compiler.New()
	if err := c.Compile(program); err != nil {
		fmt.Errorf("Compilation failed:\n %s\n", err)
		os.Exit(1)
	}

	machine := vm.New(c.Bytecode())
	if err := machine.Run(); err != nil {
		fmt.Errorf("Executing bytecode failed:\n %s\n", err)
		os.Exit(1)
	}
}
