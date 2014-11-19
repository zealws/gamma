package sexpr

import (
	"fmt"
)

/*
Type Map provides an interface for representing key-value maps as SExprs.

The Map type is designed with non-side-affecting maps in mind, but also allows for side-affecting maps to be implemented in the future.

Certain map implementations may restrict the valid key types, but the Map interface makes no such restrictions.
*/
type Map interface {
	SExpr

	/*
		Get the value of key within the map.

		If the value is found, it will return `value`, `true`.
		If the value is not found, it will return `Null`, `false`.
	*/
	Get(key SExpr) (SExpr, bool)

	/*
		Set a key to equal a value within the map.

		Returns a map with the value of `key` set to `value`.

		(The returned map may be the current map.)
	*/
	Set(key, value SExpr) Map
}

/*
Type Environ is an implementation of Map that is non-side-affecting, but is linear time for accesses.

More than one value of a key may be provided. Newer values shadow old values, but the old values are not removed.
*/
type Environ struct {
	Value SExpr
}

func NewEnviron() *Environ {
	return &Environ{Null}
}

/*
This seems like an awkward way to create an environment. Figure something else out in the future.
*/
func MakeEnviron(expressions ...SExpr) *Environ {
	if len(expressions)%2 != 0 {
		panic("MakeEnviron expects even number of parameters")
	}
	e := NewEnviron()
	for i := 0; i < len(expressions); i += 2 {
		e = e.Put(expressions[i], expressions[i+1])
	}
	return e
}

func (e *Environ) Get(key SExpr) (SExpr, bool) {
	curr := e.Value
	for {
		if IsNull(curr) {
			return Null, false
		}
		if IsEq(Caar(curr), key) {
			return Cdar(curr), true
		}
		curr = Cdr(curr)
	}
}

func (e *Environ) Set(key, value SExpr) Map {
	return e.Put(key, value)
}

func (e *Environ) Put(key, value SExpr) *Environ {
	return &Environ{Cons(Cons(key, value), e.Value)}
}

func (e *Environ) String() string {
	return fmt.Sprintf("<environ>")
}
