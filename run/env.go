package run

import (
	"axo/ast"
	"log"
)

// runtime reprentations of lambdas
type Clos struct {
        env   *local
        param string
        body  ast.Tree
}

func (c Clos) String() string {
	return "{λ...}"
}

// Globals
type global struct {
        defined ast.Defs
        evaled  ast.Defs
}

func glob_look(ge *global, k string) ast.Tree {
        ve, is_e := ge.evaled[k]
        if is_e {
                return ve
        }
        dv, is_d := ge.defined[k]
        if is_d {
                var le *local // nil
                evaluated := eval(ge, le, dv)

                // think of this next thing as caching,
                // but also marking the graph as processed
                ge.evaled[k] = evaluated
                return evaluated
        }
        log.Fatal("unbound and undefined ", k)
        return nil
}


// Locals
type local struct {
        key string
        val ast.Tree
        nxt *local
}

func extend(e *local, k string, t ast.Tree) *local {
        return &local{k, t, e}
}

func lookup(e *local, k string) ast.Tree {
        for r := e; r != nil; r = r.nxt {
                if k == r.key {
                        return r.val
                }
        }
        return nil // the empty tree nil means not found
}

// Patterns
// Pat Envs need to enforce linearity,
// ie, \x cannot be introduced twice within one pattern
// but \x \y is fine

type patEnv map[string]ast.Tree

func mergePatEnvs(l, r patEnv) patEnv {
        for k, v := range l {
                r[k] = v
        }
        return r
}

func extendPatEnv(pe patEnv, key string, val ast.Tree) {
        _, defd := pe[key]
       	if defd {
        	log.Fatal("pattern var defined in this block")
        }
        pe[key] = val
}
