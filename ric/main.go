package main

import (
	"fmt"
	"os"
)

var errors int
func Error(l Line, msg string, args ...interface{}) {
	s := fmt.Sprintf(msg, args...)
	fmt.Printf("%s:%d %s\n", l.FileName, l.LineNo, s)
	errors++
}
func Failed() bool {
	return errors > 0
}

func main() {
	prog := AstModule(Line{"<stdin>", 1}, Parse("<stdin>", os.Stdin))
	prog.Eval(globals, nil)
	for _, m := range modules {
		m.ToVerilog(os.Stdout)
	}
}
