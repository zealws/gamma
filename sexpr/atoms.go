package sexpr

import "fmt"

/*
Atom Types

An atom type is any type that self-evaluates and is eq? to itself.

That is, a value `v` is an atom if and only if:

- `(eq? v v)`
- `v` is a fixed point of the interpreter. (i.e. `(evaluate v)` yields `v`)
*/

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
	return ok && b == otherb
}

func (b Boolean) String() string {
	if b {
		return "#t"
	} else {
		return "#f"
	}
}

/**
*** Integers
**/

type Integer int64

func (i Integer) IsEq(other Comparable) bool {
	otheri, ok := other.(Integer)
	return ok && i == otheri
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", int64(i))
}

/**
*** Float
**/

type Float float64

func (f Float) IsEq(other Comparable) bool {
	otherf, ok := other.(Float)
	return ok && f == otherf
}

func (f Float) String() string {
	return fmt.Sprintf("%f", float64(f))
}
