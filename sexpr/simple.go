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
