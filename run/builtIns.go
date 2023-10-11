package run

import (
	"axo/ast"
	"fmt"
	"log"
)


// we need to wrap these up in case of unnapplied or partially applied fns
var builtInLams = map[string]ast.Tree {
        "cons": ast.ParseExpr("λf s.(cons f s)"),
        "car":  ast.ParseExpr("λl.(car l)"),
        "cdr":  ast.ParseExpr("λl.(cdr l)"),
	"+":    ast.ParseExpr("λx y.(+ x y)"),
	"*":    ast.ParseExpr("λx y.(* x y)"),
	"-":    ast.ParseExpr("λx y.(- x y)"),
	"and":  ast.ParseExpr("λx y.(and x y)"),
	"or":   ast.ParseExpr("λx y.(or x y)"),
	"not":  ast.ParseExpr("λx.(not x)"),
}

// Primitive Binary Functions

type binaryFn func(ast.Tree, ast.Tree) ast.Tree

var builtInBinary = map[string] binaryFn {
        "cons": primCons,
	"+":    primAdd,
	"*":    primMul,
	"-":    primSub,
	"and":  primAnd,
	"or":   primOr,
}

// lists are proper, with an allowance for improper lists
// of pvars and wildcards, since these just intend to describe
// what should be a proper list
func primCons(e1, e2 ast.Tree) (t ast.Tree) {
	switch e2.(type) {
	case ast.Nil, ast.Pair, ast.Wildcard, ast.PatVar:
		t = ast.Pair{e1, e2}
	default:
		fmt.Println(e2)
		log.Fatal("improper lists are not supported")
	}
	return t
}

func primAdd(e1, e2 ast.Tree) ast.Tree {
	return ast.Number(uint64(e1.(ast.Number)) + uint64(e2.(ast.Number)))
}

func primMul(e1, e2 ast.Tree) ast.Tree {
	return ast.Number(uint64(e1.(ast.Number)) * uint64(e2.(ast.Number)))
}

func primSub(e1, e2 ast.Tree) ast.Tree {
	n1 := uint64(e1.(ast.Number))
	n2 := uint64(e2.(ast.Number))
	if n2 > n1 {
		log.Fatal("subtraction error on nats")
	}
	return ast.Number(n1 - n2)
}

func primAnd(e1, e2 ast.Tree) ast.Tree {
	return ast.Boolean(bool(e1.(ast.Boolean)) && bool(e2.(ast.Boolean)) )
}
func primOr(e1, e2 ast.Tree) ast.Tree {
	return ast.Boolean(bool(e1.(ast.Boolean)) || bool(e2.(ast.Boolean)) )
}

// Primitive Unary Functions

type unaryFn func(e1 ast.Tree) ast.Tree

var builtInUnary = map[string]unaryFn {
        "car": getCar,
        "cdr": getCdr,
	"not": primNot,
}

func getCar(e1 ast.Tree) ast.Tree {
        return e1.(ast.Pair).Car
}
func getCdr(e1 ast.Tree) ast.Tree {
        return e1.(ast.Pair).Cdr
}

func primNot(e1 ast.Tree) ast.Tree {
	return ast.Boolean(!bool(e1.(ast.Boolean)))
}
