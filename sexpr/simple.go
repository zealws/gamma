package sexpr

import (
	"fmt"
)

// Primitive Non-Atom Types

/**
*** Invariant
**/

type Invariant string

func (b Invariant) IsEq(other Comparable) bool {
	otherb, ok := other.(Invariant)
	if !ok {
		return false
	}
	return b == otherb
}

func (b Invariant) String() string {
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
