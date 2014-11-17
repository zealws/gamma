package sexpr

import (
	"fmt"
)

// SExpr types made of other SExprs

/**
*** Pair
**/

type Pair struct {
	Car, Cdr SExpr
}

func Cons(car, cdr SExpr) SExpr {
	return &Pair{car, cdr}
}

func List(exprs ...SExpr) SExpr {
	result := &Pair{}
	p := result
	for i, expr := range exprs {
		if i != 0 {
			new_p := &Pair{}
			p.Cdr = new_p
			p = new_p
		}
		p.Car = expr
	}
	p.Cdr = Null
	return result
}

func (p *Pair) String() string {
	return "(" + p.privString() + ")"
}

func (p *Pair) privString() string {
	if IsNull(p.Cdr) {
		return p.Car.String()
	} else if IsPair(p.Cdr) {
		return fmt.Sprintf("%v %s", p.Car, p.Cdr.(*Pair).privString())
	}
	return fmt.Sprintf("%v . %v", p.Car, p.Cdr)
}

/**
*** Quote
**/

type QuotedExpr struct {
	Expr SExpr
}

func Quote(expr SExpr) SExpr {
	return QuotedExpr{expr}
}

func (q QuotedExpr) IsEq(expr Comparable) bool {
	qexpr, ok := expr.(QuotedExpr)
	if !ok {
		return false
	}
	return IsEq(q.Expr, qexpr.Expr)
}

func (q QuotedExpr) String() string {
	return fmt.Sprintf("'%v", q.Expr)
}
