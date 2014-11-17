package sexpr_test

import (
	. "github.com/zfjagann/gamma/sexpr"

	"testing"
)

type eqTestCase struct {
	first, second SExpr
	equals        bool
}

var eqTestCases = []eqTestCase{
	{
		List(Symbol("a"), Symbol("b")),
		List(Symbol("a"), Symbol("b")),
		false,
	},
	{
		List(Symbol("a"), Symbol("b")),
		List(Symbol("a"), Symbol("a")),
		false,
	},
	{
		Null, Null, true,
	},
	{
		Symbol("a"), Symbol("b"), false,
	},
	{
		Symbol("a"), Symbol("a"), true,
	},
	{
		True, True, true,
	},
	{
		False, False, true,
	},
	{
		True, False, false,
	},
}

func TestIsEq(t *testing.T) {
	for _, c := range eqTestCases {
		if IsEq(c.first, c.second) != c.equals {
			if c.equals {
				t.Errorf("%v != %v", c.first, c.second)
			} else {
				t.Errorf("%v == %v", c.first, c.second)
			}
		}
	}
}

type eqStarTestCase struct {
	first, second SExpr
	equals        bool
}

var eqStarTestCases = []eqStarTestCase{
	{
		List(Symbol("a"), Symbol("b")),
		List(Symbol("a"), Symbol("b")),
		true,
	},
	{
		List(Symbol("a"), Symbol("b")),
		List(Symbol("a"), Symbol("a")),
		false,
	},
	{
		Null, Null, true,
	},
	{
		Symbol("a"), Symbol("b"), false,
	},
	{
		Symbol("a"), Symbol("a"), true,
	},
	{
		True, True, true,
	},
}

func TestIsEqStar(t *testing.T) {
	for _, c := range eqStarTestCases {
		if IsEqStar(c.first, c.second) != c.equals {
			if c.equals {
				t.Errorf("%v != %v", c.first, c.second)
			} else {
				t.Errorf("%v == %v", c.first, c.second)
			}
		}
	}
}
