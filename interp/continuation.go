package interp

import (
	"fmt"
	. "github.com/zfjagann/gamma/sexpr"
)

type interpContinuation struct {
	// These arguments are named after their arguments in the original scheme interpreter.
	// They can probably be collapsed. No continuation uses more than 2 of the SExprs.
	id       string
	Expr     SExpr
	ExprList SExpr
	Answer   SExpr
	Clauses  SExpr
	Env      SExpr
	SymList  SExpr
	RandList SExpr
	Symbol   SExpr
	C        SExpr
	Rator    *Closure
}

func (c interpContinuation) String() string {
	s := "<" + c.id
	for _, val := range []SExpr{c.Expr, c.ExprList, c.Answer, c.Clauses} {
		if val != nil {
			s += fmt.Sprintf(" %v", val)
		}
	}
	if c.Rator != nil {
		s += fmt.Sprintf(" %v", c.Rator)
	}
	for _, val := range []SExpr{c.Env, c.SymList, c.RandList, c.Symbol, c.C} {
		if val != nil {
			s += fmt.Sprintf(" %v", val)
		}
	}
	s += ">"
	return s
}

// The continuation at the start of continuation
// Equivalent to (lambda (x) x)
var CID = interpContinuation{id: "cid"}

// C1 is the recursive call during a function application called after the rator has been evaluated
// C1 evaluates the parameter list, then calls C2 which calls performs the function call
func NewC1(expr, env, C SExpr) SExpr {
	return interpContinuation{id: "c1", C: C, Expr: expr, Env: env}
}

// C2 is the continuation from a function application called after the rator and randList have been evaluated
// C2 applies the function rator to the parameter list `randList`
func NewC2(answer, env, C SExpr) SExpr {
	return interpContinuation{id: "c2", C: C, Answer: answer, Env: env}
}

// C3 is the continuation from the recursize case of exprListValue
func NewC3(exprList, env, C SExpr) SExpr {
	return interpContinuation{id: "c3", C: C, ExprList: exprList, Env: env}
}

// C4 is the continuation from the recursive case of exprListValue which performs the
// cons on the result of the previous two computations
func NewC4(answer, C SExpr) SExpr {
	return interpContinuation{id: "c4", C: C, Answer: answer}
}

// C5 is the continuation from the recurisve case of condValue
func NewC5(clauses, env, C SExpr) SExpr {
	return interpContinuation{id: "c5", C: C, Clauses: clauses, Env: env}
}

// C6 is the continuation called during a closure evaluation with the argument list
func NewC6(rator *Closure, C SExpr) SExpr {
	return interpContinuation{id: "c6", C: C, Rator: rator}
}

// C7 is the continuation from the recursive case of augmentedEnv
func NewC7(symList, randList, C SExpr) SExpr {
	return interpContinuation{id: "c7", C: C, SymList: symList, RandList: randList}
}

// C8 is called during a define block with the evaluated expression
func NewC8(symbol, C SExpr) SExpr {
	return interpContinuation{id: "c8", C: C, Symbol: symbol}
}
