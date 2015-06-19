package sexpr

type Comparable interface {
	SExpr
	IsEq(other Comparable) bool
}

type SExpr interface {
	String() string
}
