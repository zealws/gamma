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
		List(Cons(Symbol("foo"), Symbol("bar")))),
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
		mustParse("(- 1 1)"),
		Integer(0)),
	pass(
		mustParse("(- 4 2 1)"),
		Integer(1)),

	/**
	*** Negative Test Cases
	**/
	fail(
		Symbol("a"),
		`environment lookup failed for symbol "a"`),
	fail(
		List(Symbol("car"), List(Symbol("a"), Symbol("b"))),
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
		`missing expr in cons: (cons a)`),
	fail(
		mustParse("(cons)"),
		`missing expr in cons: (cons)`),
	fail(
		mustParse("(eq? 'a)"),
		`missing expr in eq?: (eq? a)`),
	fail(
		mustParse("(eq?)"),
		`missing expr in eq?: (eq?)`),
	fail(
		mustParse("(symbol?)"),
		`missing expr in symbol?: (symbol?)`),
	fail(
		mustParse("(null?)"),
		`missing expr in null?: (null?)`),
	fail(
		mustParse("(apply 'a)"),
		`missing expr in apply: (apply a)`),
	fail(
		mustParse("(apply)"),
		`missing expr in apply: (apply)`),
	fail(
		mustParse("(call/cc)"),
		`missing expr in call/cc: (call/cc)`),
}

func TestDefinesSymbol(t *testing.T) {
	interp := NewInterpreter(DefaultEnvironment)
	_, err := interp.Evaluate(mustParse("(define a '(a))"))
	if err != nil {
		t.Fatal(err)
	}
	expr, err := interp.Evaluate(mustParse("a"))
	if err != nil {
		t.Fatal(err)
	}
	expected := List(Symbol("a"))
	if !IsEqStar(expr, expected) {
		t.Fatalf("Expected %v but was %v", expected, expr)
	}
}

/**
*** Utilities Below This Point
**/

func mustParse(input string) SExpr {
	expr, err := parse.Parse(input)
	if err != nil {
		panic(fmt.Sprintf("Could not parse test input %q: %v", input, err))
	}
	return expr
}

type testCase struct {
	env   SExpr
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
				return fmt.Sprintf("Expected %v but was %v", expected, actual)
			} else if !IsEqStar(expected, actual) {
				return fmt.Sprintf("Expected %v but was %v", expected, actual)
			}
			return ""
		},
	}
}

func passEnv(input SExpr, expected SExpr, env SExpr) testCase {
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

func failEnv(input SExpr, msg string, env SExpr) testCase {
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
