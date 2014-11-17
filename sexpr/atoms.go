package sexpr

// Atom Types

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
