package ast

import (
        "fmt"
        "log"
        "os"
	"strconv"
)

func parseDefs(ds Defs, ts tokens) Defs {
        for ts.peek().typeOf != eof {
                assert(ts, define)
                name := assert(ts, iden).literal
                expr := parseExpr(ts, false) // non pattern top
                _, prev_defd := ds[name]
                if prev_defd {
			fmt.Println(name)
                        log.Fatal("duplicate definition")
                }
                ds[name] = expr
        }
        assert(ts, eof) // check clean exit
        return ds
}

func parseExpr(ts tokens, inPatt bool) Tree {
        var t Tree
        switch ts.peek().typeOf {
        case iden:
                t = Iden(ts.take().literal)
        case symbol:
                t = Symbol(ts.take().literal)
	case number:
		ns := ts.take().literal
		n, _ := strconv.ParseUint(ns, 10, 64)
		t = Number(n)
	case booltrue:
		ts.take()
		t = Boolean(true)
	case boolfalse:
		ts.take()
		t = Boolean(false)
        case lpar:
                t = parseApps(ts, inPatt)
        case lambda:
                ts.take()
                t = parseLambdas(ts, inPatt)
	case let:
		t = parseLet(ts, inPatt)
	case ifexpr:
		ts.take()
		t = parseIfExpr(ts, inPatt)
        case lbrack:
                ts.take()
                t = parseList(ts, inPatt)
        case match:
                t = parseMatch(ts, inPatt)
	case backslash:
		check(inPatt, ts.take())
		t = PatVar(assert(ts, iden).literal)
	case wildcard:
		check(inPatt, ts.take())
		t = Wildcard{}
        default:
                log.Fatal("parse error! ", ts.peek())
        }
        return t
}


func parseLambdas(ts tokens, inPatt bool) Tree {
        if ts.peek().typeOf == dot {
                ts.take()
                return parseExpr(ts, inPatt)
        }
        idn := assert(ts, iden).literal
        return Lambda{idn, parseLambdas(ts, inPatt)}
}

func parseApps(ts tokens, inPatt bool) Tree {
        var t Tree
        ts.take()
        l := parseExpr(ts, inPatt)
        r := parseExpr(ts, inPatt)
        t = Apply{l, r}
        for ts.peek().typeOf != rpar {
                t = Apply{t, parseExpr(ts, inPatt)}
        }
        ts.take()
        return t
}

func parseList(ts tokens, inPatt bool) Tree {
        if ts.peek().typeOf == rbrack {
                ts.take()
                return Nil{}
        }
        el := parseExpr(ts, inPatt)
        return Apply{Apply{Iden("cons"), el}, parseList(ts, inPatt)}
}

// let is overloaded, idens cause a basic let binding,
// otherwise we assume destructured assignment, which 
// we just convert to a pattern match on one case,
// keeps the syntax small without any additional cases in eval

func parseLet(ts tokens, inPatt bool) Tree {
	line_no := ts.take().line_no
	var t Tree
	if ts.peek().typeOf == iden {
           name := ts.take().literal
           b := parseExpr(ts, inPatt)
           i := parseExpr(ts, inPatt)
           t = Let{name, b, i}
	} else {
	   p := parseExpr(ts, true) // a pattern case
	   v := parseExpr(ts, inPatt)
	   i := parseExpr(ts, inPatt)
	   t = Match{v, []Tree{p}, []Tree{i}, line_no}
	}
	return t
}

func parseIfExpr(ts tokens, inPatt bool) Tree {
	t := parseExpr(ts, inPatt)
	c := parseExpr(ts, inPatt)
	a := parseExpr(ts, inPatt)
	return IfExpr{t, c, a}
}

func parseMatch(ts tokens, inPatt bool) Tree {
        line_no := ts.take().line_no
        matchOn := parseExpr(ts, inPatt)
        var tests []Tree
        var thens []Tree
        for {
                if ts.peek().typeOf == pipe {
                        ts.take()
                        tests = append(tests, parseExpr(ts, true)) // patt case
                        thens = append(thens, parseExpr(ts, inPatt))
                } else {
                        break
                }
        }
        if len(tests) == 0 {
                log.Fatal("no tests given in match statement line whatever")
        }
        return Match{matchOn, tests, thens, line_no}
}

// helpers
func check(inPatt bool, t token) {
        if !inPatt {
                msg := "token %v used outside match expression line %v\n"
                fmt.Printf(msg, t.typeOf, t.line_no)
                os.Exit(1)
        }
}

func assert(ts tokens, toktype tokType) token {
        tok := ts.take()
        if tok.typeOf != toktype {
                msg := "parse fail line %v\nwanted %v, got %v\n"
                fmt.Printf(msg, tok.line_no, toktype, tok.typeOf)
                os.Exit(1)
        }
        return tok
}
