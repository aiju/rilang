//go:generate goyacc -o parse.tab.go parse.y

package main

import (
	"bufio"
	"unicode"
	"unicode/utf8"
	"math/big"
	"regexp"
)

type Line struct {
	FileName string
	LineNo int
}

type lex struct {
	r *bufio.Reader
}

var kwtab = map[string] int {
	"module": TMODULE,
	"int": TINT,
	"if": TIF,
	"else": TELSE,
}

var optab = map[string] int {
	"==": TEQ,
}
func makeCharTab(tab map[string] int) (r map[rune] bool) {
	r = make(map[rune] bool)
	for s, _ := range tab {
		rune, _ := utf8.DecodeRuneInString(s)
		r[rune] = true
	}
	return
}
var opchar = makeCharTab(optab)

var curline Line

func isIdent(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c >= 0x80 || c == '_'
}

var uN_regexp, _ = regexp.Compile("^u[0-9]+$")

func (l *lex) Lex(lval *yySymType) int {
	var c rune
	var err error

	r := l.r
	for {
		c, _, err = r.ReadRune()
		if err != nil { return -1 }
		if c == '\n' { curline.LineNo++ }
		if !unicode.IsSpace(c) { break }
	}
	if c >= '0' && c <= '9' {
		s := string(c)
		for {
			c, _, err = r.ReadRune()
			if err != nil { break }
			if !isIdent(c) { break }
			if c != '_' { s += string(c) }
		}
		r.UnreadRune()
		lval.num = big.NewInt(0)
		_, ok := lval.num.SetString(s, 0)
		if !ok { l.Error(err.Error()) }
		return TNUM
	}
	if isIdent(c) {
		s := string(c)
		for {
			c, _, err = r.ReadRune()
			if err != nil { break }
			if !isIdent(c) { break }
			s += string(c)
		}
		r.UnreadRune()
		if tok, ok := kwtab[s]; ok {
			return tok
		}
		if s[0] == 'u' && uN_regexp.MatchString(s) {
			lval.num = big.NewInt(0)
			lval.num.SetString(s[1:], 0)
			return TUX
		}
		lval.str = s
		return TSYM
	}
	if opchar[c] {
		d, _, err := r.ReadRune()
		if err != nil { return int(c) }
		if tok, ok := optab[string(c) + string(d)]; ok {
			return tok
		}
		r.UnreadRune()
	}
	return int(c)
}

func (l *lex) Error(s string) {
	Error(curline, "%s", s)
}
