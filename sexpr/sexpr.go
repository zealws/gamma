package sexpr

import (
	"fmt"
)

var (
	SimpleString bool = false
)

type Comparable interface {
	IsEq(other Comparable) bool
}

/*
Implementations:
	Null (value)
	Symbol
	Boolean
	*Pair
	*Builtin
*/
type SExpr interface {
	// IsEq(other SExpr) bool
	String() string
}

/**
*** Null
**/

type nullt struct{}

var (
	Null SExpr = nullt{}
)

func (nullt) IsEq(other Comparable) bool {
	_, ok := other.(nullt)
	return ok
}

func (nullt) String() string {
	return "<null>"
}

/**
*** Symbol
**/

type Symbol string

func (sym Symbol) IsEq(other Comparable) bool {
	othersym, ok := other.(Symbol)
	if !ok {
		return false
	}
	return sym == othersym
}

func (sym Symbol) String() string {
	return string(sym)
}

/**
*** Boolean
**/

type Boolean bool

var (
	True         SExpr = Boolean(true)
	TrueLiteral  SExpr = True
	False        SExpr = Boolean(false)
	FalseLiteral SExpr = False
)

func (b Boolean) IsEq(other Comparable) bool {
	otherb, ok := other.(Boolean)
	if !ok {
		return false
	}
	return b == otherb
}

func (b Boolean) String() string {
	if b {
		return "#t"
	} else {
		return "#f"
	}
}

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
	if SimpleString {
		return fmt.Sprintf("(%v . %v)", p.Car, p.Cdr)
	}
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
*** Interpreter Builtin
**/

type Builtin string

func (b Builtin) IsEq(other Comparable) bool {
	otherb, ok := other.(Builtin)
	if !ok {
		return false
	}
	return b == otherb
}

func (b Builtin) String() string {
	return fmt.Sprintf("<built-in %s>", string(b))
}

/**
*** Continuation
**/

type Continuation struct {
	C SExpr
}

func NewContinuation(C SExpr) Continuation {
	return Continuation{C}
}

func (b Continuation) String() string {
	return fmt.Sprintf("<cont %v>", b.C)
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
