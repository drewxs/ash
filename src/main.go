package main

import (
	"ash/repl"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Ash 0.0.1")
	repl.Start(os.Stdin, os.Stdout)
}
