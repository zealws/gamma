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

		Returns a map with all the values in the original map, as well as
		the value of `key` set to `value`.

		(The returned map may be the same object as the original map, but
		need not necessarily be to allow for one-way-mutable maps.)
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
Builds an environment from the given map. The keys in the map are converted to
symbols before being added to the environment.
*/
func BuildSymbolEnviron(exprs map[string]SExpr) *Environ {
	e := NewEnviron()
	for k, v := range exprs {
		e = e.Put(Symbol(k), v)
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
