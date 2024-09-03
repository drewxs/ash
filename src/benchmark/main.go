package main

import (
	"flag"
	"fmt"
	"time"

	"ash/compiler"
	"ash/evaluator"
	"ash/lexer"
	"ash/object"
	"ash/parser"
	"ash/vm"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

var input = `
let fib = fn(n) {
    if n == 0 {
        0
    } else {
        if n == 1 {
            1
        } else {
            fib(n - 1) + fib(n - 2)
        }
    }
};
fib(35)
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		vm := vm.New(comp.Bytecode())

		start := time.Now()

		err = vm.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = vm.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, result=%s, duration=%s\n",
		*engine, result.Inspect(), duration)
}
