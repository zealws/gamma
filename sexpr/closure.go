package sexpr

import (
	"fmt"
)

type Closure struct {
	SymList SExpr
	Body    SExpr
	Env     SExpr
}

func NewClosure(symList, body, env SExpr) *Closure {
	return &Closure{symList, body, env}
}

func (c *Closure) String() string {
	return fmt.Sprintf("<closure %v %v %v>", c.SymList, c.Body, c.Env)
}
