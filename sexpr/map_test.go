package sexpr

import (
	"testing"
)

func TestEnviron(t *testing.T) {
	var _ Map = NewEnviron()

	e := NewEnviron()
	e = e.Put(Symbol("a"), Symbol("x"))
	e = e.Put(Symbol("b"), Symbol("y"))
	e = e.Put(Symbol("c"), Symbol("z"))

	assertGetEq(t, e, Symbol("a"), Symbol("x"))
	assertGetEq(t, e, Symbol("b"), Symbol("y"))
	assertGetEq(t, e, Symbol("c"), Symbol("z"))

	e = MakeEnviron(
		Symbol("a"), Symbol("x"),
		Symbol("b"), Symbol("y"),
		Symbol("c"), Symbol("z"))

	assertGetEq(t, e, Symbol("a"), Symbol("x"))
	assertGetEq(t, e, Symbol("b"), Symbol("y"))
	assertGetEq(t, e, Symbol("c"), Symbol("z"))
}

func assertGetEq(t *testing.T, m Map, key, exp SExpr) {
	act, ok := m.Get(key)
	if !ok {
		t.Fatal("Could not get %v from map", key)
	}
	if !IsEq(act, exp) {
		t.Fatal("Expected %v but got %v", act, exp)
	}
}
