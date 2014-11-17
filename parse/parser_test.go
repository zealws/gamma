package parse_test

import (
	"github.com/zfjagann/gamma/parse"
	. "github.com/zfjagann/gamma/sexpr"
	"testing"
)

var parseTestCases = []parseTestCase{
	Pass("'abc", Quote(Symbol("abc"))),
	Pass("1234", Integer(1234)),
	Pass("0x123", Integer(291)),
	Pass("(a b c)", List(Symbol("a"), Symbol("b"), Symbol("c"))),

	Fail("'", "unexpected EOF in symbol expression at offset 2"),
	Fail("(a b", "unexpected EOF in list at offset 6"),
	Fail("(a b .)", "unexpected ')'. expecting symbol at offset 9"),
}

/**
*** Test Utilities
**/

type parseTestCase struct {
	input  string
	output SExpr
	err    string
}

func (c parseTestCase) Do(t *testing.T) {
	expr, err := parse.Parse(c.input)
	if err != nil {
		if c.err == "" {
			t.Errorf("Could not parse %q: %v", c.input, err)
		} else if err.Error() != c.err {
			t.Errorf("Expected %q but got %q for %q", c.err, err.Error(), c.input)
		}
	} else if c.err != "" {
		t.Errorf("Expected failure %q but got %q for %q", c.err, c.output, c.input)
	} else if expr == nil && expr != c.output {
		t.Errorf("Expected %v but got %v for %q", c.output, expr, c.input)
	} else if expr != nil && !IsEqStar(expr, c.output) {
		t.Errorf("Expected %v but got %v for %q", c.output, expr, c.input)
	}
}

func Pass(input string, output SExpr) parseTestCase {
	return parseTestCase{input, output, ""}
}

func Fail(input string, err string) parseTestCase {
	return parseTestCase{input, nil, err}
}

func TestParsesSymbol(t *testing.T) {
	for _, c := range parseTestCases {
		c.Do(t)
	}
}
