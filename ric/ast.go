//go:generate stringer -type=AstKind
//go:generate stringer -type=Oper
package main

import (
	"math/big"
	"fmt"
)

type AstKind int
const (
	ASTNIL AstKind = iota
	ASTSYM
	ASTNUM
	ASTUX
	ASTBIN
	ASTMODULE
	ASTASS
	ASTDEF
)
type Oper int
const (
	OPNIL Oper = iota
	OPEQ
	OPADD
	OPSUB
	OPMUL
	OPDIV
	OPMOD
	OPASS
)
type Ast struct {
	Line
	Kind AstKind
	Name string
	Num *big.Int
	N []*Ast
}
func AstSym(l Line, s string) *Ast { return &Ast{Line: l, Kind: ASTSYM, Name: s} }
func AstNum(l Line, n *big.Int) *Ast { return &Ast{Line: l, Kind: ASTNUM, Num: n} }
func AstUx(l Line, n *big.Int) *Ast { return &Ast{Line: l, Kind: ASTUX, Num: n} }
func AstModule(l Line, n []*Ast) *Ast { return &Ast{Line: l, Kind: ASTMODULE, N: n} }
func AstAss(l Line, lhs *Ast, rhs *Ast) *Ast { return &Ast{Line: l, Kind: ASTASS, N: []*Ast{lhs, rhs} } }
func AstDef(l Line, name string, typ *Ast, attr []*Ast) *Ast { return &Ast{Line: l, Kind: ASTDEF, Name: name, N: append([]*Ast{ typ }, attr...) } }

type operator struct {
	name string
	prec int
	right bool
}
var opertab = map[Oper]operator{
	OPEQ: {"==", 100, false},
	OPADD: {"+", 200, false},
	OPSUB: {"-", 200, false},
	OPMUL: {"*", 210, false},
	OPDIV: {"/", 210, false},
	OPMOD: {"%", 210, false},
	OPASS: {"=", 50, true},
}
func (o Oper) tab() operator {
	t, ok := opertab[o]
	if ok {
		return t
	}
	return operator{
		o.String(), 0, false,
	}
}

func parens(s string, opprec int, env int) string {
	if opprec < env {
		return "(" + s + ")"
	}
	return s
}

func tabs(ind int) string {
	s := ""
	for i := 0; i < ind; i++ {
		s += "\t"
	}
	return s
}

func (p *Ast) PrettyPrint(terse bool, prec int, ind int) string {
	if p == nil { return "<nil>" }
	switch p.Kind {
	case ASTNIL:
		return "nil"
	case ASTSYM:
		return p.Name
	case ASTNUM:
		return p.Num.String()
	case ASTUX:
		return fmt.Sprintf("u%v", p.Num)
	case ASTASS:
		s := p.N[0].PrettyPrint(terse, 51, ind) + " = " + p.N[1].PrettyPrint(terse, 50, ind)
		return parens(s, 50, prec)
	case ASTDEF:
		s := p.Name + " :"
		if p.N[0] != nil { s += " " + p.N[0].PrettyPrint(terse, 451, ind) }
		return parens(s, 450, prec)
	case ASTMODULE:
		if terse {
			return "module {...}"
		}
		s := "module {\n"
		for _, n := range p.N {
			s += tabs(ind+1) + n.PrettyPrint(false, 0, ind+1) + ";\n"
		}
		s += tabs(ind) + "}"
		return s
	default:
		return p.Kind.String()
	}
}

func (p *Ast) String() string {
	return p.PrettyPrint(true, 0, 0)
}

func (p *Ast) Walk(pre func(*Ast) *Ast, post func(*Ast) *Ast) *Ast {
	if p == nil { return nil }
	if pre != nil { p = pre(p) }
	q := &*p
	keep := false
	q.N = nil
	for i := range q.N {
		e := p.N[i].Walk(pre, post)
		if e != p.N[i] { keep = true }
		if e != nil { q.N = append(q.N, e) }
	}
	if keep { p = q }
	if post != nil { p = post(p) }
	return p
}
