package utils

import "io"

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, Format(RED, "error: "))
		io.WriteString(out, msg+"\n")
	}
}
