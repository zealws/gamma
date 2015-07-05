package interp

import (
	"fmt"
	"time"

	. "github.com/zfjagann/gamma/sexpr"
)

var (
	lambdaLiteral SExpr = Symbol("lambda")
	defineLiteral SExpr = Symbol("define")
	condLiteral   SExpr = Symbol("cond")
	elseLiteral   SExpr = Symbol("else")
	ifLiteral     SExpr = Symbol("if")

	DefaultEnvironment *Environ = BuildSymbolEnviron(map[string]SExpr{
		"car":     Invariant("car"),
		"cdr":     Invariant("cdr"),
		"cons":    Invariant("cons"),
		"eq?":     Invariant("eq?"),
		"symbol?": Invariant("symbol?"),
		"null?":   Invariant("null?"),
		"apply":   Invariant("apply"),
		"call/cc": Invariant("call/cc"),
		"exit":    Invariant("exit"),
		"env":     Invariant("env"),
		"time":    Invariant("time"),
		"sleep":   Invariant("sleep"),
		"+":       builtin{"+", Sum},
		"-":       builtin{"-", Subtract},
		"*":       builtin{"*", Product},
		"/":       builtin{"/", Quotient},
		//"^": Invariant("^"),
		//"%": Invariant("%"),
	})

	Exit error = fmt.Errorf("interpreter exited")
)

type Interpreter struct {
	env *Environ
}

func NewInterpreter(env *Environ) *Interpreter {
	return &Interpreter{env}
}

func (in *Interpreter) define(symbol, expr SExpr) SExpr {
	in.env = in.env.Put(symbol, expr)
	return Null
}

func (in *Interpreter) Evaluate(expr SExpr) (SExpr, error) {
	return in.schemeValue(in.env, newInterpStack(), expr)
}

func (in *Interpreter) schemeValue(env *Environ, stack *interpStack, expr SExpr) (result SExpr, err error) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("panic: %v\n%s\n", e, getStack())
			result = nil
			err = fmt.Errorf("panic: %v", e)
		}
	}()
	var (
		clauses, exprList, C, randList, rator, sym, symList, answer SExpr
		found                                                       bool
		answerEnv                                                   *Environ
	)

	C = CID
	goto exprValue

exprValue:
	// evaluate the expression `expr` with regard to `env` and call `C` with the result
	stack.trace("exprValue(expr,env,C)", expr, env, C)

	if IsAtom(expr) {
		// Atoms are fixed-points of the interpreter
		answer = expr
		goto applyC
	} else if IsSymbol(expr) {
		sym = expr
		goto symValue
	} else if q, ok := expr.(QuotedExpr); ok {
		answer = q.Expr
		goto applyC
	} else if IsEq(Car(expr), condLiteral) {
		clauses = Cdr(expr)
		goto condValue
	} else if IsEq(Car(expr), ifLiteral) {
		_len := randLength(Cdr(expr))
		if _len < 3 {
			return nil, fmt.Errorf("missing parameter from if statement: %v", expr)
		} else if _len > 3 {
			return nil, fmt.Errorf("extra parameters from if statement: %v", expr)
		}
		_cond, err := ECadr(expr)
		if err != nil {
			return nil, fmt.Errorf("CCCC: %v", expr)
		}
		_exprs, err := ECddr(expr)
		if err != nil {
			return nil, fmt.Errorf("AAAA: %v", expr)
		}
		C = NewC9(_exprs, C)
		expr = _cond
		goto exprValue
	} else if IsEq(Car(expr), lambdaLiteral) {
		argList, err := ECadr(expr)
		if err != nil {
			return nil, fmt.Errorf("missing parameter list in function literal: %v", expr)
		}
		body, err := ECaddr(expr)
		if err != nil {
			return nil, fmt.Errorf("missing body in function literal: %v", expr)
		}
		answer = NewClosure(argList, body, env)
		goto applyC
	} else if IsEq(Car(expr), defineLiteral) {
		defSym, err := ECadr(expr)
		if err != nil {
			return nil, fmt.Errorf("missing symbol in define: %v", expr)
		}
		defExpr, err := ECaddr(expr)
		if err != nil {
			return nil, fmt.Errorf("missing expression in define: %v", expr)
		}
		C = NewC8(defSym, C)
		expr = defExpr
		goto exprValue
	} else {
		C = NewC1(expr, env, C)
		expr = Car(expr)
		goto exprValue
	}

exprListValue:
	// evaluate all the expressions in `exprList` and call `C` with the result
	stack.trace("exprListValue(exprList,C)", exprList, C)

	if IsEq(exprList, Null) {
		answer = Null
		goto applyC
	} else {
		C = NewC3(exprList, env, C)
		expr = Car(exprList)
		goto exprValue
	}

symValue:
	// perform an environment lookup of `sym` within `env` and call `C` with the result
	stack.trace("symValue(sym,env,C)", sym, env, C)

	answer, found = env.Get(sym)
	if !found {
		return nil, fmt.Errorf("environment lookup failed for symbol %q", sym.(Symbol))
	} else {
		goto applyC
	}

condValue:
	// evaluate the cond block defined by `clauses` and call `C` with the result
	stack.trace("condValue(clauses,C)", clauses, C)

	if IsNull(clauses) {
		// (cond)
		return nil, fmt.Errorf("invalid empty cond block")
	} else if IsNull(Car(clauses)) {
		// (cond ())
		return nil, fmt.Errorf("invalid empty cond condition")
	}
	{
		clause := Car(clauses)
		condition, err := ECar(clause)
		if err != nil {
			// (cond ())
			return nil, fmt.Errorf("missing condition in cond clause: %v", clause)
		}
		condExpr, err := ECadr(clause)
		if err != nil {
			// (cond (x))
			return nil, fmt.Errorf("missing expression in cond clause: %v", clause)
		}
		if IsEq(condition, elseLiteral) {
			expr = condExpr
			goto exprValue
		} else {
			C = NewC5(clauses, env, C)
			expr = condition
			goto exprValue
		}
	}

appValue:
	// apply the operation `rator` with `randList` as arguments and call `C` with the result
	stack.trace("appValue(rator,randList,C)", rator, randList, C)

	if bi, ok := rator.(builtin); ok {
		exprList = randList
		C = bi.cont(C)
		goto exprListValue
	} else if bi, ok := rator.(Invariant); ok {
		switch string(bi) {
		case "car":
			if err := checkLen(1, rator, randList); err != nil {
				return nil, err
			}
			answer, err = ECaar(randList)
			if err != nil {
				return nil, err
			}
			goto applyC
		case "cdr":
			if err := checkLen(1, rator, randList); err != nil {
				return nil, err
			}
			answer, err = ECdar(randList)
			if err != nil {
				return nil, err
			}
			goto applyC
		case "cons":
			if err := checkLen(2, rator, randList); err != nil {
				return nil, err
			}
			f := Car(randList)
			s := Cadr(randList)
			answer = Cons(f, s)
			goto applyC
		case "eq?":
			if err := checkLen(2, rator, randList); err != nil {
				return nil, err
			}
			f := Car(randList)
			s := Cadr(randList)
			answer = IsEqExpr(f, s)
			goto applyC
		case "symbol?":
			if err := checkLen(1, rator, randList); err != nil {
				return nil, err
			}
			f := Car(randList)
			answer = IsSymbolExpr(f)
			goto applyC
		case "null?":
			if err := checkLen(1, rator, randList); err != nil {
				return nil, err
			}
			f := Car(randList)
			answer = IsNullExpr(f)
			goto applyC
		case "apply":
			if err := checkLen(2, rator, randList); err != nil {
				return nil, err
			}
			rator = Car(randList)
			randList = Cadr(randList)
			goto appValue
		case "env":
			if err := checkLen(0, rator, randList); err != nil {
				return nil, err
			}
			answer = env
			goto applyC
		case "exit":
			if err := checkLen(0, rator, randList); err != nil {
				return nil, err
			}
			return nil, Exit
		case "call/cc":
			if err := checkLen(1, rator, randList); err != nil {
				return nil, err
			}
			rator = Car(randList)
			randList = List(NewContinuation(C))
			goto appValue
		case "time":
			if err := checkLen(0, rator, randList); err != nil {
				return nil, err
			}
			answer = Integer(time.Now().UnixNano() / 1000000)
			goto applyC
		case "sleep":
			if err := checkLen(1, rator, randList); err != nil {
				return nil, err
			}
			f := Car(randList)
			if t, ok := f.(Integer); ok {
				time.Sleep(time.Duration(t) * time.Second)
			} else {
				return nil, fmt.Errorf("Invalid time value: %v", f)
			}
			answer = Integer(time.Now().UnixNano() / 1000000)
			goto applyC
		default:
			return nil, fmt.Errorf("unknown built-in method: %q", string(bi))
		}
	} else if clos, ok := rator.(*Closure); ok {
		C = NewC6(clos, C)
		symList = clos.SymList
		env = clos.Env
		goto augmentedEnv
	} else if cont, ok := rator.(Continuation); ok {
		C = cont.C
		answer = Car(randList)
		goto applyC
	} else {
		return nil, fmt.Errorf("Unknown operator: %v", rator)
	}

augmentedEnv:
	// augment the environment `env` with symbols `symList` and values `randList` and call `C` with the modified environment
	stack.trace("augmentedEnv(symList,randList,env,C)", symList, randList, env, C)

	answerEnv = env
	for {
		if IsNull(symList) {
			goto applyC
		} else if IsSymbol(symList) {
			answerEnv = answerEnv.Put(symList, randList)
			symList = Null
			randList = Null
		} else {
			if err := checkLen(randLength(symList), rator, randList); err != nil {
				return nil, err
			}
			answerEnv = answerEnv.Put(Car(symList), Car(randList))
			symList = Cdr(symList)
			randList = Cdr(randList)
		}
	}

applyC:
	// apply the continuation `C` to the value `answer`
	// the continuation values defined below are a result of continuation function literals
	// in the original scheme interpreter that have been translated to Go
	stack.trace("applyC(answer,C)", answer, C)

	if C == CID {
		// the program has finished computing
		return answer, nil
	}
	if c, ok := C.(*builtinContinuation); ok {
		answer, err = c.b.f(answer)
		if err != nil {
			return nil, err
		}
		C = c.C
		goto applyC
	} else if c, ok := C.(interpContinuation); ok {
		switch c.id {
		case "c1":
			// C1 is the recursive call during a function application called after the rator has been evaluated
			// C1 evaluates the parameter list, then calls C2 which calls performs the function call
			exprList = Cdr(c.Expr)
			env = c.Env
			C = NewC2(answer, c.Env, c.C)
			goto exprListValue
		case "c2":
			// C2 is the continuation from a function application called after the rator and randList have been evaluated
			// C2 applies the function rator to the parameter list `randList`
			rator = c.Answer
			randList = answer
			env = c.Env
			C = c.C
			goto appValue
		case "c3":
			// C3 is the continuation from the recursive case of exprListValue which performs the tail recursion
			exprList = Cdr(c.ExprList)
			env = c.Env
			C = NewC4(answer, c.C)
			goto exprListValue
		case "c4":
			// C4 is the continuation from the recursive case of exprListValue which performs the
			// cons on the result of the previous two computations
			answer = Cons(c.Answer, answer)
			C = c.C
			goto applyC
		case "c5":
			// C5 is the continuation from the recursive case of condValue
			if !IsEq(answer, False) {
				expr = Cadar(c.Clauses)
				env = c.Env
				C = c.C
				goto exprValue
			} else {
				clauses = Cdr(c.Clauses)
				env = c.Env
				C = c.C
				goto condValue
			}
		case "c6":
			// C6 is the continuation called during a closure evaluation with the environment
			expr = c.Rator.Body
			env = answerEnv
			C = c.C
			goto exprValue
		case "c8":
			// C8 is called during a define block with the evaluated expression
			if clos, ok := answer.(*Closure); ok {
				// Cheap hack to make recursive functions work
				clos.Env = clos.Env.Put(c.Symbol, clos)
			}
			in.define(c.Symbol, answer)
			answer = Null
			C = c.C
			goto applyC
		case "c9":
			// C9 is called during an if with the evaluated condition
			if IsEq(answer, False) {
				expr = Cadr(c.ExprList)
			} else {
				expr = Car(c.ExprList)
			}
			C = c.C
			goto exprValue
		default:
			return nil, fmt.Errorf("invalid continuation value: %v", C)
		}
	} else {
		return nil, fmt.Errorf("invalid continuation value: %v", C)
	}
}

func randLength(randList SExpr) int {
	l := 0
	for {
		if IsNull(randList) {
			return l
		}
		l += 1
		randList = Cdr(randList)
	}
}

func checkLen(size int, rator, randList SExpr) error {
	actual := randLength(randList)
	if actual != size {
		return fmt.Errorf("%v expects %d arguments but was given %d", rator, size, actual)
	}
	return nil
}
