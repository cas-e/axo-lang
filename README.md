# The Axo Programming Language

Axo is a **toy** functional programming language, built in Go, created as a learning exercise. 

Some of its features include:

* first-class lambda functions with support for currying
* pattern matching
* destructured assignment style let-bindings
* unicode support for source code 

And it supports the basic types:

* numbers, symbols, lists, booleans, and closures

Some major missing features that may be a to-do someday:

* currently missing static types, only has run-time type checking
* no string type at all
* has a slow, tree walking interpreter. a compiled-to-VM version could be next


# Running

Currently, running `go run main.go` will run the `sudoko.axo` script, which is a simple sudoku puzzle solver written in Axo. Tested with Go 1.18. 