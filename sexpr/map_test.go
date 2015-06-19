package sexpr

import (
	"testing"
)

func TestEnviron(t *testing.T) {
	// Type check.
	var _ Map = NewEnviron()

	e := NewEnviron()
	e = e.Put(Symbol("a"), Symbol("x"))
	e = e.Put(Symbol("b"), Symbol("y"))
	e = e.Put(Symbol("c"), Symbol("z"))

	assertGetEq(t, e, Symbol("a"), Symbol("x"))
	assertGetEq(t, e, Symbol("b"), Symbol("y"))
	assertGetEq(t, e, Symbol("c"), Symbol("z"))

	e = MakeEnviron(
		Symbol("a"), Symbol("j"),
		Symbol("b"), Symbol("k"),
		Symbol("c"), Symbol("l"))

	assertGetEq(t, e, Symbol("a"), Symbol("j"))
	assertGetEq(t, e, Symbol("b"), Symbol("k"))
	assertGetEq(t, e, Symbol("c"), Symbol("l"))
}

func assertGetEq(t *testing.T, m Map, key, exp SExpr) {
	act, ok := m.Get(key)
	if !ok {
		t.Fatalf("Could not get %v from map", key)
	}
	if !IsEq(act, exp) {
		t.Fatalf("Expected %v but got %v", exp, act)
	}
}
