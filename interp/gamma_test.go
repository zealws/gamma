package interp

import (
	"fmt"
	"github.com/zfjagann/gamma/parse"
	. "github.com/zfjagann/gamma/sexpr"
	"testing"
)

const TestTraceSize = 2 // How many stack records to show on failure

var testCases = []testCase{
	/**
	*** Positive Test Cases
	**/
	passEnv(
		Symbol("foo"),
		Symbol("bar"),
		MakeEnviron(Symbol("foo"), Symbol("bar"))),
	pass(
		mustParse("'(a b c d)"),
		List(Symbol("a"), Symbol("b"), Symbol("c"), Symbol("d"))),
	pass(
		mustParse("(car '(a b))"),
		Symbol("a")),
	pass(
		mustParse("(cdr '(a b))"),
		List(Symbol("b"))),
	pass(
		mustParse("(cons 'a '(b))"),
		List(Symbol("a"), Symbol("b"))),
	pass(
		mustParse("(cons 'a 'b)"),
		Cons(Symbol("a"), Symbol("b"))),
	pass(
		mustParse("(call/cc (lambda (c) (c 'a)))"),
		Symbol("a")),
	pass(
		mustParse("(call/cc (lambda (cc) ((lambda (y) (cc 'bar)) 'foo)))"),
		Symbol("bar")),
	pass(
		mustParse("(null? '())"),
		True),
	pass(
		mustParse("(eq? 'a 'a)"),
		True),
	pass(
		mustParse("(eq? 'a 'b)"),
		False),
	pass(
		mustParse("(symbol? 'a)"),
		True),
	pass(
		mustParse("(symbol? '())"),
		False),
	pass(
		mustParse("(null? '())"),
		True),
	pass(
		mustParse("(null? 'a)"),
		False),
	pass(
		mustParse("((lambda (x) (cond (x 'a) (else 'b))) #t)"),
		Symbol("a")),
	pass(
		mustParse("((lambda (x) (cond (x 'a) (else 'b))) #f)"),
		Symbol("b")),
	pass(
		mustParse("(+ 1 1)"),
		Integer(2)),
	pass(
		mustParse("(+ 1 1 15)"),
		Integer(17)),
	pass(
		mustParse("(* 2 5)"),
		Integer(10)),
	pass(
		mustParse("(* 2 5 3)"),
		Integer(30)),
	pass(
		mustParse("(/ 16 4)"),
		Integer(4)),
	pass(
		mustParse("(/ 36 4 3)"),
		Integer(3)),
	pass(
		mustParse("(- 1 1)"),
		Integer(0)),
	pass(
		mustParse("(- 4 2 1)"),
		Integer(1)),
	passEnv(
		mustParse("(env)"),
		MakeEnviron(Symbol("env"), Invariant("env"), Symbol("foo"), Symbol("bar")),
		MakeEnviron(Symbol("env"), Invariant("env"), Symbol("foo"), Symbol("bar"))),
	pass(
		mustParse("((lambda x 'foo) 'bar)"),
		Symbol("foo")),
	pass(
		mustParse("((lambda x x) 'foo 'bar 'baz)"),
		List(Symbol("foo"), Symbol("bar"), Symbol("baz"))),

	/**
	*** Negative Test Cases
	**/
	fail(
		Symbol("a"),
		`environment lookup failed for symbol "a"`),
	fail(
		mustParse("(car (a b))"),
		`environment lookup failed for symbol "a"`),
	fail(
		mustParse("(lambda)"),
		`missing parameter list in function literal: (lambda)`),
	fail(
		mustParse("(lambda a)"),
		`missing body in function literal: (lambda a)`),
	fail(
		mustParse("(define)"),
		`missing symbol in define: (define)`),
	fail(
		mustParse("(define a)"),
		`missing expression in define: (define a)`),
	fail(
		mustParse("(cond)"),
		`invalid empty cond block`),
	fail(
		mustParse("(cond (else))"),
		`missing expression in cond clause: (else)`),
	fail(
		mustParse("(cond (x))"),
		`missing expression in cond clause: (x)`),
	fail(
		mustParse("(cond ())"),
		`invalid empty cond condition`),
	fail(
		mustParse("(cons 'a)"),
		`<built-in cons> expects 2 arguments but was given 1`),
	fail(
		mustParse("(cons)"),
		`<built-in cons> expects 2 arguments but was given 0`),
	fail(
		mustParse("(eq? 'a)"),
		`<built-in eq?> expects 2 arguments but was given 1`),
	fail(
		mustParse("(eq?)"),
		`<built-in eq?> expects 2 arguments but was given 0`),
	fail(
		mustParse("(symbol?)"),
		`<built-in symbol?> expects 1 arguments but was given 0`),
	fail(
		mustParse("(null?)"),
		`<built-in null?> expects 1 arguments but was given 0`),
	fail(
		mustParse("(apply 'a)"),
		`<built-in apply> expects 2 arguments but was given 1`),
	fail(
		mustParse("(apply)"),
		`<built-in apply> expects 2 arguments but was given 0`),
	fail(
		mustParse("(call/cc)"),
		`<built-in call/cc> expects 1 arguments but was given 0`),
	fail(
		mustParse("((lambda (x) 'a))"),
		`<closure> expects 1 arguments but was given 0`),
}

func TestDefinesSymbol(t *testing.T) {
	interp := NewInterpreter(DefaultEnvironment)
	assertEvaluates(t, interp, "(define a '(a))", nil)
	assertEvaluates(t, interp, "a", List(Symbol("a")))
}

func TestDefinesRecursiveFunction(t *testing.T) {
	interp := NewInterpreter(DefaultEnvironment)
	assertEvaluates(t, interp, "(define len (lambda (x) (cond ((null? x) 0) (else (+ 1 (len (cdr x)))))))", nil)
	assertEvaluates(t, interp, "(len '(a b c d))", Integer(4))
}

func TestCanFormatRecursiveFunction(t *testing.T) {
	interp := NewInterpreter(DefaultEnvironment)
	assertEvaluates(t, interp, "(define len (lambda (x) (cond ((null? x) 0) (else (+ 1 (len (cdr x)))))))", nil)
	expr := assertEvaluates(t, interp, "len", nil)
	// If this call recurses, bad things happen.
	expr.String()
}

/**
*** Utilities Below This Point
**/

func assertEvaluates(t *testing.T, interp *Interpreter, input string, expected SExpr) SExpr {
	expr, err := interp.Evaluate(mustParse(input))
	if err != nil {
		t.Fatal(err)
	}
	if expected != nil && !IsEqStar(expr, expected) {
		t.Fatalf("Expected %v but was %v", expected, expr)
	}
	return expr
}

func mustParse(input string) SExpr {
	expr, err := parse.Parse(input)
	if err != nil {
		panic(fmt.Sprintf("Could not parse test input %q: %v", input, err))
	}
	return expr
}

type testCase struct {
	env   *Environ
	check func(*testing.T, SExpr, error) string
	input SExpr
}

func pass(input SExpr, expected SExpr) testCase {
	return testCase{
		env:   DefaultEnvironment,
		input: input,
		check: func(t *testing.T, actual SExpr, err error) string {
			if err != nil {
				return fmt.Sprintf("Could not evaluate %v: %v", input, err)
			} else if actual == nil {
				return fmt.Sprintf("Expected %v but was %v for %v", expected, actual, input)
			} else if !IsEqStar(expected, actual) {
				return fmt.Sprintf("Expected %v but was %v for %v", expected, actual, input)
			}
			return ""
		},
	}
}

func passEnv(input, expected SExpr, env *Environ) testCase {
	return testCase{
		env:   env,
		input: input,
		check: func(t *testing.T, actual SExpr, err error) string {
			if err != nil {
				return fmt.Sprintf("Could not evaluate %v: %v", input, err)
			} else if actual == nil {
				return fmt.Sprintf("Expected %v but was %v", expected, actual)
			} else if !IsEqStar(expected, actual) {
				return fmt.Sprintf("Expected %v but was %v", expected, actual)
			}
			return ""
		},
	}
}

func fail(input SExpr, msg string) testCase {
	return testCase{
		env:   DefaultEnvironment,
		input: input,
		check: func(t *testing.T, actual SExpr, err error) string {
			if err == nil {
				return fmt.Sprintf("Expected %q but was %v", msg, err)
			} else if err.Error() != msg {
				return fmt.Sprintf("Expected %q but was %q", msg, err.Error())
			}
			return ""
		},
	}
}

func failEnv(input SExpr, msg string, env *Environ) testCase {
	return testCase{
		env:   env,
		input: input,
		check: func(t *testing.T, actual SExpr, err error) string {
			if err == nil {
				return fmt.Sprintf("Expected %q but was %v", msg, err)
			} else if err.Error() != msg {
				return fmt.Sprintf("Expected %q but was %q", msg, err.Error())
			}
			return ""
		},
	}
}

func TestEvaluatesTestCases(t *testing.T) {
	for _, c := range testCases {
		interp := NewInterpreter(c.env)
		expr, err := interp.Evaluate(c.input)
		fail := c.check(t, expr, err)
		if fail != "" {
			fmt.Printf("--- TRACE(%d) ---\n%v\n--- END TRACE ---\n", TestTraceSize, interp.Trace().Last(TestTraceSize))
			t.Error(fail)
		}
	}
}
