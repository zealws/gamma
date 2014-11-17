package sexpr

type Comparable interface {
	IsEq(other Comparable) bool
}

type SExpr interface {
	String() string
}
