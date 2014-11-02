package sexpr

import (
	"fmt"
)

/**
*** Equality Comparisons
**/

func IsEq(a, b SExpr) bool {
	ca, ok := a.(Comparable)
	if !ok {
		return false
	}
	cb, ok := b.(Comparable)
	if !ok {
		return false
	}
	return ca.IsEq(cb)
}

func IsEqStar(rawa, rawb SExpr) bool {
	if pa, ok := rawa.(*Pair); ok {
		if pb, ok := rawb.(*Pair); ok {
			return IsEqStar(pa.Car, pb.Car) && IsEqStar(pa.Cdr, pb.Cdr)
		}
	}
	return IsEq(rawa, rawb)
}

/**
*** Type Checks
**/

func IsSymbol(e SExpr) bool {
	_, ok := e.(Symbol)
	return ok
}

func IsNull(e SExpr) bool {
	return e == Null
}

func IsPair(e SExpr) bool {
	_, ok := e.(*Pair)
	return ok
}

/**
*** Expression Wrappers
**/

func IsEqExpr(a, b SExpr) SExpr {
	return Boolean(IsEq(a, b))
}

func IsSymbolExpr(e SExpr) SExpr {
	return Boolean(IsSymbol(e))
}

func IsNullExpr(e SExpr) SExpr {
	return Boolean(IsNull(e))
}

func IsPairExpr(e SExpr) SExpr {
	return Boolean(IsPair(e))
}

/**
*** Car/Cdr Variants
**/

func ECar(e SExpr) (SExpr, error) {
	p, ok := e.(*Pair)
	if !ok {
		return nil, fmt.Errorf("car on non-pair: %v", e)
	}
	return p.Car, nil
}

func ECdr(e SExpr) (SExpr, error) {
	p, ok := e.(*Pair)
	if !ok {
		return nil, fmt.Errorf("cdr on non-pair: %v", e)
	}
	return p.Cdr, nil
}

func ECadr(e SExpr) (SExpr, error) {
	p, err := ECdr(e)
	if err != nil {
		return nil, err
	}
	return ECar(p)
}

func ECdar(e SExpr) (SExpr, error) {
	p, err := ECar(e)
	if err != nil {
		return nil, err
	}
	return ECdr(p)
}

func ECaar(e SExpr) (SExpr, error) {
	p, err := ECar(e)
	if err != nil {
		return nil, err
	}
	return ECar(p)
}

func ECaddr(e SExpr) (SExpr, error) {
	p, err := ECdr(e)
	if err != nil {
		return nil, err
	}
	p, err = ECdr(p)
	if err != nil {
		return nil, err
	}
	return ECar(p)
}

func ECadar(e SExpr) (SExpr, error) {
	p, err := ECar(e)
	if err != nil {
		return nil, err
	}
	p, err = ECdr(p)
	if err != nil {
		return nil, err
	}
	return ECar(p)
}

// Only use if you are absolutely sure that `e` is a pair
func Car(e SExpr) SExpr {
	r, err := ECar(e)
	if err != nil {
		panic(err)
	}
	return r
}

// Only use if you are absolutely sure that `e` is a pair
func Cdr(e SExpr) SExpr {
	r, err := ECdr(e)
	if err != nil {
		panic(err)
	}
	return r
}

// Only use if you are absolutely sure that `e` is a pair
func Caar(e SExpr) SExpr {
	return Car(Car(e))
}

// Only use if you are absolutely sure that `e` is a pair
func Cdar(e SExpr) SExpr {
	return Cdr(Car(e))
}

// Only use if you are absolutely sure that `e` is a pair
func Cadar(e SExpr) SExpr {
	return Car(Cdr(Car(e)))
}
