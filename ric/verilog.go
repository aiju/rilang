package main

import (
	"io"
	"fmt"
)

func (n *Node) ToVerilog(w io.Writer) {
	switch n.Kind {
	case NODESYM:
		fmt.Fprintf(w, "%s", n.Symbol.Name)
	default:
		panic("Node.ToVerilog: unknown " + n.Kind.String())
	}
}

func (m *Module) ToVerilog(w io.Writer) {
	fmt.Fprint(w, "module " + m.Name + "(\n")
	for _, s := range m.Signals {
		fmt.Fprint(w, s.Name + "\n")
	}
	fmt.Fprint(w, ");\n")
	for _, n := range m.N {
		switch n.Kind {
		case NODEASS:
			fmt.Fprint(w, "\tassign ")
			n.N[0].ToVerilog(w)
			fmt.Fprint(w, " = ")
			n.N[1].ToVerilog(w)
			fmt.Fprint(w, ";\n")
		}
	}
	fmt.Fprint(w, "endmodule\n\n")
}
