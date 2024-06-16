package repl

import (
	"ash/compiler"
	"ash/lexer"
	"ash/object"
	"ash/parser"
	color "ash/utils"
	"ash/vm"
	"bufio"
	"fmt"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Printf(color.Format(color.CYAN, ">> "))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		vm := vm.NewWithGlobalsStore(code, globals)
		err = vm.Run()
		if err != nil {
			fmt.Fprintf(out, "Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := vm.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, color.Format(color.RED, "error: "))
		io.WriteString(out, msg+"\n")
	}
}
