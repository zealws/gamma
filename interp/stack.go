package interp

import (
	"fmt"
	"github.com/zfjagann/gamma/sexpr"
	"runtime"
	"strings"
)

const TraceMaxSize int = 1024 // 1K max frames in a stack trace

/**
*** Interpreter methods for stack traces
**/

func (in *Interpreter) trace(spec string, args ...sexpr.SExpr) {
	if in.nextTrace {
		in.stack.Enqueue(&StackFrame{spec, args})
	}
	in.nextTrace = true
}

func (in *Interpreter) ignoreNextTrace() {
	in.nextTrace = false
}

func (in *Interpreter) Trace() StackTrace {
	trace := make(StackTrace, 0, TraceMaxSize)
	for _, v := range in.stack.Values() {
		trace = append(trace, v.(*StackFrame))
	}
	return trace
}

/**
*** Utilities
**/

type StackFrame struct {
	spec  string
	exprs []sexpr.SExpr
}

func (f *StackFrame) String() string {
	var (
		name    string   = f.spec
		columns []string = nil
	)
	idx := strings.Index(f.spec, "(")
	if idx != -1 {
		name = f.spec[0:idx]
		columns = strings.Split(f.spec[idx+1:len(f.spec)-1], ",")
	}
	str := "\033[32m" + name + "\033[39m"
	for i, v := range f.exprs {
		str += "\n  "
		if columns != nil {
			str += columns[i] + " => "
		}
		str += fmt.Sprintf("%v", v)
	}
	return str
}

type StackTrace []*StackFrame

func (t StackTrace) Last(n int) StackTrace {
	start := len(t) - 1 - n
	end := len(t) - 1
	if start < 0 {
		start = 0
	}
	return StackTrace(t[start:end])
}

func (t StackTrace) String() string {
	s := ""
	for i, f := range t {
		if i != 0 {
			s += "\n"
		}
		s += f.String()
	}
	return s
}

func getStack() string {
	buf := make([]byte, 1024, 1024)
	i := runtime.Stack(buf, false)
	return string(buf[0:i])
}
