package interp

import (
	. "github.com/zfjagann/gamma/sexpr"

	"fmt"
	"github.com/zfjagann/golang-ring"
)

var (
	lambdaLiteral SExpr = Symbol("lambda")
	defineLiteral SExpr = Symbol("define")
	condLiteral   SExpr = Symbol("cond")
	elseLiteral   SExpr = Symbol("else")

	DefaultEnvironment SExpr = List(
		Cons(Symbol("car"), Invariant("car")),
		Cons(Symbol("cdr"), Invariant("cdr")),
		Cons(Symbol("cons"), Invariant("cons")),
		Cons(Symbol("eq?"), Invariant("eq?")),
		Cons(Symbol("symbol?"), Invariant("symbol?")),
		Cons(Symbol("null?"), Invariant("null?")),
		Cons(Symbol("apply"), Invariant("apply")),
		Cons(Symbol("call/cc"), Invariant("call/cc")),
		Cons(Symbol("exit"), Invariant("exit")),
	)

	Exit error = fmt.Errorf("interpreter exited")
)

type Interpreter struct {
	env       SExpr
	stack     *ring.Ring
	nextTrace bool
}

func NewInterpreter(env SExpr) *Interpreter {
	stack := &ring.Ring{}
	stack.SetCapacity(TraceMaxSize)
	return &Interpreter{env, stack, true}
}

func (in *Interpreter) define(symbol, expr SExpr) SExpr {
	in.env = Cons(Cons(symbol, expr), in.env)
	return Null
}

func (in *Interpreter) Evaluate(expr SExpr) (SExpr, error) {
	return in.schemeValue(in.env, expr)
}

func (in *Interpreter) schemeValue(env, expr SExpr) (result SExpr, err error) {
	defer func() {
		e := recover()
		if e != nil {
			result = nil
			err = fmt.Errorf("panic: %v", e)
		}
	}()
	var (
		clauses, exprList, C, randList, rator, sym, symList, answer SExpr
	)

	C = CID
	goto exprValue

exprValue:
	// evaluate the expression `expr` with regard to `env` and call `C` with the result
	in.trace("exprValue(expr,env,C)", expr, env, C)

	if IsEq(expr, TrueLiteral) {
		answer = TrueLiteral
		goto applyC
	} else if IsEq(expr, FalseLiteral) {
		answer = FalseLiteral
		goto applyC
	} else if IsNull(expr) {
		answer = Null
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
	in.trace("exprListValue(exprList,C)", exprList, C)

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
	in.trace("symValue(sym,env,C)", sym, env, C)

	if IsEq(env, Null) {
		return nil, fmt.Errorf("environment lookup failed for symbol %q", sym.(Symbol))
	} else if IsEq(Caar(env), sym) {
		answer = Cdar(env)
		goto applyC
	} else {
		env = Cdr(env)
		in.ignoreNextTrace()
		goto symValue
	}

condValue:
	// evaluate the cond block defined by `clauses` and call `C` with the result
	in.trace("condValue(clauses,C)", clauses, C)

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
	in.trace("appValue(rator,randList,C)", rator, randList, C)

	if bi, ok := rator.(Invariant); ok {
		switch string(bi) {
		case "car":
			answer, err = ECaar(randList)
			if err != nil {
				return nil, err
			}
			goto applyC
		case "cdr":
			answer, err = ECdar(randList)
			if err != nil {
				return nil, err
			}
			goto applyC
		case "cons":
			f, err := ECar(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			s, err := ECadr(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			answer = Cons(f, s)
			goto applyC
		case "eq?":
			f, err := ECar(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			s, err := ECadr(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			answer = IsEqExpr(f, s)
			goto applyC
		case "symbol?":
			f, err := ECar(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			answer = IsSymbolExpr(f)
			goto applyC
		case "null?":
			f, err := ECar(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			answer = IsNullExpr(f)
			goto applyC
		case "apply":
			f, err := ECar(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			s, err := ECadr(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			rator = f
			randList = s
			goto appValue
		case "exit":
			return nil, Exit
		case "call/cc":
			f, err := ECar(randList)
			if err != nil {
				return nil, fmt.Errorf("missing expr in %s: %v", string(bi), Cons(Symbol(string(bi)), randList))
			}
			rator = f
			randList = List(NewContinuation(C))
			goto appValue
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
	in.trace("augmentedEnv(symList,randList,env,C)", symList, randList, env, C)

	if IsEq(symList, Null) {
		answer = env
		goto applyC
	} else {
		C = NewC7(symList, randList, C)
		symList = Cdr(symList)
		randList = Cdr(randList)
		goto augmentedEnv
	}

applyC:
	// call the continuation `C` to the value `answer`
	// the continuation values defined below are a result of continuation function literals
	// in the original scheme interpreter that have been transliterated to Go
	in.trace("applyC(answer,C)", answer, C)

	if C == CID {
		// the program has finished computing
		return answer, nil
	}
	c, ok := C.(interpContinuation)
	if !ok {
		return nil, fmt.Errorf("invalid continuation value: %v", C)
	}
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
		if !IsEq(answer, FalseLiteral) {
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
		// C6 is the continuation called during a closure evaluation with the argument list
		expr = c.Rator.Body
		env = answer
		C = c.C
		goto exprValue
	case "c7":
		// C7 is the continuation from the recursive case of augmentedEnv
		answer = Cons(Cons(Car(c.SymList), Car(c.RandList)), answer)
		C = c.C
		goto applyC
	case "c8":
		// C8 is called during a define block with the evaluated expression
		in.define(c.Symbol, answer)
		answer = Null
		C = c.C
		goto applyC
	default:
		return nil, fmt.Errorf("invalid continuation value: %v", C)
	}
}
