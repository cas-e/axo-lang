package run

import (
        "axo/ast"
        "fmt"
        "log"
)

func EvalProg(ds ast.Defs) ast.Tree {
        main_expr, defd := ds["main"]
        if !defd {
                log.Fatal("main is not defined in file")
        }
        gl_evaled := make(ast.Defs)
        gl_env := &global{defined: ds, evaled: gl_evaled}

        var lc *local // nil
        main_evaled := eval(gl_env, lc, main_expr)
        return main_evaled
}

func eval(ge *global, le *local, expr ast.Tree) ast.Tree {
        var t ast.Tree
        switch e := expr.(type) {

        case ast.Iden:
		idn := string(e)
                idv := lookup(le, idn)
                if idv == nil {
                        lam, isLam := builtInLams[idn]
                        if isLam {
                               t = eval(ge, le, lam)
                        } else {
                               t = glob_look(ge, idn)
                        }
                } else {
                        t = idv
                }


	// the parser enforces the fact that PatVar is never used
	// in the wrong context
	case ast.Symbol, ast.Nil, ast.PatVar, ast.Wildcard, ast.Number, ast.Boolean:
		t = e // immeadiate values

        case ast.Lambda:
                t = Clos{le, e.Var, e.Body}

        case ast.Apply:
                unfn, isUn := checkUnaryPrim(e.Rator)
                binfn, isBin := checkAllBinArgs(e.Rator)

                if isUn {
                        arg := eval(ge, le, e.Rand)
                        t = unfn(arg)
                } else if isBin {
                        e1 := eval(ge, le, e.Rator.(ast.Apply).Rand)
                        e2 := eval(ge, le, e.Rand)
                        t = binfn(e1, e2)

                } else {
                        rtr := eval(ge, le, e.Rator)
                        rnd := eval(ge, le, e.Rand)
                        t = apply(ge, rtr, rnd)
                }

	case ast.Let:
		letval := eval(ge, le, e.Be)
		t = eval(ge, extend(le, e.Name, letval), e.In)

	case ast.IfExpr:
		testRes := bool(eval(ge, le, e.Test).(ast.Boolean))
		if testRes {
			t = eval(ge, le, e.Consq)
		} else {
			t = eval(ge, le, e.Altern)
		}

	case ast.Match:
		matchOn := eval(ge, le, e.MatchOn)
		var isMatch bool
		for i, v := range e.Tests {
			patt := eval(ge, le, v)
			pe := make(patEnv)
			isMatch, pe = checkPattern(pe, patt, matchOn)
			if isMatch {
				for k, v := range pe {
					le = extend(le, k, v)
				}
				t = eval(ge, le, e.Thens[i])
				break
			}
		}
		
		if !isMatch {
			fmt.Println("missed pattern")
			fmt.Println("defined match on line ", e.Line)
			fmt.Println(matchOn)
			log.Fatal("pgrm term")
		}
        default:
                log.Fatal("how did i get here in eval")


        }
        return t
}

func apply(g *global, clos ast.Tree, arg ast.Tree) ast.Tree {
        var t ast.Tree
        switch c := clos.(type) {
        case Clos:
                t = eval(g, extend(c.env, c.param, arg), c.body)
        default:
                fmt.Println(clos)
                fmt.Println("is not a function")

                log.Fatal("how did i get here in apply")
        }
        return t
}

func checkPattern(pe patEnv, patt, expr ast.Tree) (bool, patEnv) {

	switch p := patt.(type) {
	case ast.Wildcard:
		return true, pe
	case ast.PatVar:
		extendPatEnv(pe, string(p), expr)
		return true, pe

	}

	if isClos(patt) || isClos(expr) {
		log.Fatal("closures are not supported in checkPattern")
		var empty patEnv
		return false, empty
	} 
	
	if isPair(patt) && isPair(expr) {
                lbool, lenv := checkPattern(pe, patt.(ast.Pair).Car, expr.(ast.Pair).Car)
                rbool, renv := checkPattern(pe, patt.(ast.Pair).Cdr, expr.(ast.Pair).Cdr)
                return lbool && rbool, mergePatEnvs(lenv, renv)
	}

	// symbols and nils should support equality already, and ints too ofc
	if patt == expr {
		return true, pe
	}

	empty := make(patEnv)
	return false, empty
}


// Helpers

func isPair(t ast.Tree) bool {
	_, chk := t.(ast.Pair)
	return chk
}

func isClos(t ast.Tree) bool {
	_, chk := t.(Clos)
	return chk
}

// look for (@ (@ op e1) e2)... but called from (@ ... )
func checkAllBinArgs(t ast.Tree) (binaryFn, bool) {
        a, ok := t.(ast.Apply)
        if ok {
                i, ok := a.Rator.(ast.Iden)
                if ok {
                        fn, ok := builtInBinary[string(i)]
                        return fn, ok
                }
                return nil, false
        }
        return nil, false
}

func checkUnaryPrim(t ast.Tree) (unaryFn, bool) {
        i, ok := t.(ast.Iden)
        if ok {
                fn, ok := builtInUnary[string(i)]
                return fn, ok
        }
        return nil, false
}
