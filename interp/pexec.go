package interp

import (
	. "github.com/zfjagann/gamma/sexpr"
)

func (in *Interpreter) makePexec(env *Environ, val SExpr) SExpr {
	// FIXME strong typing!
	comms := make(chan *pexecResult, 1)

	var rcomms chan<- *pexecResult = comms

	go func() {
		result := &pexecResult{}
		result.stack = newInterpStack()
		result.expr, result.err = in.schemeValue(env, result.stack, val)
		rcomms <- result
		close(rcomms)
	}()
	return &Pexec{comms, nil}
}

type Pexec struct {
	// FIXME strong typing!
	comms  <-chan *pexecResult
	result *pexecResult
}

func (*Pexec) String() string {
	return "<pexec>"
}

func (p *Pexec) GetResult() (SExpr, error) {
	if p.result == nil {
		p.result = <-p.comms
	}

	if p.result.err != nil {
		return nil, &PexecError{p.result.err, p.result.stack}
	} else {
		return p.result.expr, nil
	}
}

type pexecResult struct {
	expr  SExpr
	err   error
	stack *interpStack
}

type PexecError struct {
	Err   error
	Stack *interpStack
}

func (p *PexecError) Error() string {
	return p.Err.Error()
}
