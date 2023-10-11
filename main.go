package main

import (
	"axo/ast"
	"axo/run"
	"fmt"
)

func main() {

	tree := ast.ParseFile("sudoku.axo")

	value := run.EvalProg(tree)
	
	fmt.Println(value)

}



