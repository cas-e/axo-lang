package ast

import (
	"fmt"
	"strings"
	"strconv"
)

// exports
func ParseFile(path string) Defs {
	// get prelude
	empty := make(Defs)
	pre := parseDefs(empty, lexFile("prelude.axo"))
        return parseDefs(pre, lexFile(path))
}

func ParseExpr(code string) Tree {
        return parseExpr(lexBytes([]byte(code)), false) // outside patt
}


type Tree interface {
	String() string
}
type Defs map[string]Tree

// lambdas
type Iden string
type Lambda struct {
	Var  string
	Body Tree
}
type Apply struct {
	Rator Tree
	Rand  Tree
}

// bindings
type Let struct {
	Name string
	Be   Tree
	In   Tree
}

// data types
type Symbol string
type Number uint64
type Boolean bool

type Nil struct{}
type Pair struct{
	Car Tree
	Cdr Tree
}
// if expressions
type IfExpr struct {
	Test   Tree
	Consq  Tree
	Altern Tree
}

// pattern matching
type PatVar string
type Wildcard struct{}

type Match struct {
	MatchOn Tree
	Tests   []Tree
	Thens   []Tree
	Line    int
}

// printing
func (i Iden) String() string {
        return fmt.Sprintf("{idn %v}", string(i))
}
func (a Apply) String() string {
        return fmt.Sprintf("{@ %v %v}", a.Rator, a.Rand)
}

func (l Lambda) String() string {
        return fmt.Sprintf("{λ%v. %v}", l.Var, l.Body)
}

func (l Let) String() string {
	b := l.Be.String()
	i := l.In.String()
	return fmt.Sprintf("let %v %v\n%v\n", l.Name, b, i)
}

func (s Symbol) String() string {
	return fmt.Sprintf("%v", string(s)) 	
}
func (n Number) String() string {
	return strconv.Itoa(int(n))
}
func (b Boolean) String() string {
	var s string
	if b {
		s = "true"
	} else {
		s = "false"
	}
	return fmt.Sprintf("%v", s)
}

func (i IfExpr) String() string {
	ts := i.Test.String()
	cs := i.Consq.String()
	as := i.Altern.String()
	return fmt.Sprintf("? %v %v %v\n", ts, cs, as)
}
func (n Nil) String() string {
        return fmt.Sprintf("[]")
}


// pairs handled at end 

func (w Wildcard) String() string {
        return fmt.Sprintf("_")
}
func (p PatVar) String() string {
	return fmt.Sprintf("{pvar %v}", string(p))
}

func (m Match) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("\n{μ %v\n", m.MatchOn))
	for i, v := range m.Tests {
		s.WriteString(fmt.Sprintf("| %v ", v))
		s.WriteString(fmt.Sprintf("-> %v\n", m.Thens[i]))
	}
	return s.String()
}

// printing lists
func (p Pair) String() string {
        s := strings.Builder{}
        s.WriteString("[")
        listToString(p, &s)
        return s.String()
}

// i don't need a return value because i am just side-effecting this builder
func listToString(t Tree, s *strings.Builder) {
        l := t.(Pair).Car
        r := t.(Pair).Cdr

        s.WriteString(l.String())
	
        restToString(r, s)
}

func restToString(t Tree, s *strings.Builder) {
	p, isP := t.(Pair)
        if (t == Nil{}) {
                s.WriteString("]")
        } else if isP {
		s.WriteString(" ")
                s.WriteString(p.Car.String())
                restToString(p.Cdr, s)
        } 
}

