package sexpr

type Closure struct {
	SymList SExpr
	Body    SExpr
	Env     *Environ
}

func NewClosure(symList, body SExpr, env *Environ) *Closure {
	return &Closure{symList, body, env}
}

func (*Closure) String() string {
	return "<closure>"
}

type Thunk interface {
	SExpr
	GetResult() (SExpr, error)
}
