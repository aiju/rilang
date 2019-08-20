//go:generate stringer -type=NodeKind
package main

import (
	"math/big"
	"fmt"
)

type Symbol struct {
	Line
	Name string
	Value *Node
	SymTab *SymTab
	Scope *SymTab
}
type SymTab struct {
	Name string
	Tab map[string] *Symbol
	Up *SymTab
}

func NewSymTab(up *SymTab) *SymTab {
	return &SymTab{
		Tab: make(map[string] *Symbol),
		Up: up,
	}
}
func (st *SymTab) Lookup(name string) *Symbol {
	for t := st; t != nil; t = t.Up {
		if s, ok := t.Tab[name]; ok {
			return s
		}
	}
	return nil
}
func (s *Symbol) FullName() string {
	r := s.Name
	for st := s.Scope; st.Up != globals; st = st.Up {
		if st.Name == "" {
			r = "???." + r
		} else {
			r = st.Name + "." + r
		}
	}
	return r
}

var tempidx int
func TempName() string {
	tempidx++
	return fmt.Sprintf("__%d", tempidx)
}

func NewSymbol(l Line, name string, scope *SymTab) *Symbol {
	s := &Symbol{
		Line: l,
		Name: name,
		Scope: scope,
	}
	scope.Tab[name] = s
	return s
}

var globals = NewSymTab(nil)

type NodeKind int
const (
	NODENIL NodeKind = iota
	NODENUM
	NODEASS
	NODESYM
	NODEMOD
)
type Node struct {
	Line
	Kind NodeKind
	Name string
	Num *big.Int
	N []*Node
	Symbol *Symbol
	Module *Module
}
func NodeAss(l Line, lhs *Node, rhs *Node) *Node { return &Node{Line: l, Kind: NODEASS, N: []*Node{lhs, rhs} } }
func NodeSym(l Line, sym *Symbol) *Node { return &Node{Line: l, Kind: NODESYM, Symbol: sym} }
func NodeMod(l Line, mod *Module) *Node { return &Node{Line: l, Kind: NODEMOD, Module: mod} }
type Module struct {
	Name string
	SymTab *SymTab
	N []*Node
	Signals []*Symbol
}
var modules []*Module

func (p *Node) PrettyPrint(terse bool, prec int, ind int) string {
	if p == nil { return "<nil>" }
	switch p.Kind {
	case NODENIL:
		return "nil"
	case NODESYM:
		return p.Symbol.Name
	case NODENUM:
		return p.Num.String()
	case NODEASS:
		s := p.N[0].PrettyPrint(terse, 51, ind) + " = " + p.N[1].PrettyPrint(terse, 50, ind)
		return parens(s, 50, prec)
	case NODEMOD:
		return "module " + p.Module.Name
	default:
		return p.Kind.String()
	}
}
func (p *Node) String() string {
	return p.PrettyPrint(true, 0, 0)
}

func (m *Module) String() string {
	s := "module " + m.Name + " {\n"
	for i := range m.N {
		s += m.N[i].String() + ";\n"
	}
	return s + "}"
}

func (p *Ast) Eval(scope *SymTab, mod *Module) *Node {
	if p == nil { return nil }
	switch p.Kind {
	case ASTMODULE:
		m := &Module{SymTab: NewSymTab(scope)}
		if scope.Up == nil {
			m.Name = "__top"
		} else {
			m.Name = TempName()
		}
		m.SymTab.Name = m.Name
		for i := range p.N {
			e := p.N[i].Eval(m.SymTab, m)
			if e != nil { m.N = append(m.N, e) }
		}
		modules = append(modules, m)
		return NodeMod(p.Line, m)
	case ASTASS:
		lhs := p.N[0].Eval(scope, mod)
		rhs := p.N[1].Eval(scope, mod)
		if lhs.Kind == NODESYM && rhs.Kind == NODEMOD {
			rhs.Module.Name = lhs.Symbol.FullName()
			rhs.Module.SymTab.Name = rhs.Module.Name
			lhs.Symbol.Value = rhs
			return rhs
		}
		return NodeAss(p.Line, lhs, rhs)
	case ASTDEF:
		sym := NewSymbol(p.Line, p.Name, scope)
		mod.Signals = append(mod.Signals, sym)
		return NodeSym(p.Line, sym)
	case ASTSYM:
		s := scope.Lookup(p.Name)
		if s == nil { Error(p.Line, "%s undefined", s) }
		return NodeSym(p.Line, s)
	default:
		panic("Ast.Eval: unknown " + p.Kind.String())
	}
}
