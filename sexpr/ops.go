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
	} else if ea, ok := rawa.(*Environ); ok {
		if eb, ok := rawb.(*Environ); ok {
			return IsEqStar(ea.Value, eb.Value)
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

// See atoms.go for a definition of an atom
func IsAtom(e SExpr) bool {
	switch e.(type) {
	case Boolean:
		return true
	case nullt:
		return true
	case Integer:
		return true
	case Float:
		return true
	}
	return false
}

func IsNumber(e SExpr) bool {
	switch e.(type) {
	case Integer:
		return true
	case Float:
		return true
	}
	return false
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
func Cadr(e SExpr) SExpr {
	return Car(Cdr(e))
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

/**
*** Math Primitives
**/

// Adds the two provided numbers, returns an error if they are not of numerical types.
func Plus(a, b SExpr) (SExpr, error) {
	switch at := a.(type) {
	case Integer:
		switch bt := b.(type) {
		case Integer:
			return Integer(int64(at) + int64(bt)), nil
		case Float:
			return Float(float64(int64(at)) + float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot add type %T", b)
		}
	case Float:
		switch bt := b.(type) {
		case Integer:
			return Float(float64(at) + float64(int64(bt))), nil
		case Float:
			return Float(float64(at) + float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot add type %T", b)
		}
	default:
		return nil, fmt.Errorf("Cannot add type %T", a)
	}
}

// Sums all the numbers in the provided list. Returns an error if they are not of numerical types.
func Sum(top SExpr) (SExpr, error) {
	var err error
	var result SExpr = Integer(0)
	cur := top
	for {
		if IsNull(cur) {
			break
		}
		result, err = Plus(Car(cur), result)
		if err != nil {
			return nil, err
		}
		cur = Cdr(cur)
	}
	return result, nil
}

// Subtracts the two provided numbers, returns an error if they are not of numerical types.
func Minus(a, b SExpr) (SExpr, error) {
	switch at := a.(type) {
	case Integer:
		switch bt := b.(type) {
		case Integer:
			return Integer(int64(at) - int64(bt)), nil
		case Float:
			return Float(float64(int64(at)) - float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot subtract type %T", b)
		}
	case Float:
		switch bt := b.(type) {
		case Integer:
			return Float(float64(at) - float64(int64(bt))), nil
		case Float:
			return Float(float64(at) - float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot subtract type %T", b)
		}
	default:
		return nil, fmt.Errorf("Cannot subtract type %T", a)
	}
}

// Subtracts all the numbers in the provided list. Returns an error if they are not of numerical types.
func Subtract(top SExpr) (SExpr, error) {
	var err error
	if IsNull(top) {
		return nil, fmt.Errorf("subtraction expects at least one parameter")
	}
	var result SExpr = Car(top)
	cur := Cdr(top)
	for {
		if IsNull(cur) {
			break
		}
		result, err = Minus(result, Car(cur))
		if err != nil {
			return nil, err
		}
		cur = Cdr(cur)
	}
	return result, nil
}

// Multiplies the two provided numbers, returns an error if they are not of numerical types.
func Multiply(a, b SExpr) (SExpr, error) {
	switch at := a.(type) {
	case Integer:
		switch bt := b.(type) {
		case Integer:
			return Integer(int64(at) * int64(bt)), nil
		case Float:
			return Float(float64(int64(at)) * float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot multiply type %T", b)
		}
	case Float:
		switch bt := b.(type) {
		case Integer:
			return Float(float64(at) * float64(int64(bt))), nil
		case Float:
			return Float(float64(at) * float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot multiply type %T", b)
		}
	default:
		return nil, fmt.Errorf("Cannot multiply type %T", a)
	}
}

// Multiplies all the numbers in the provided list. Returns an error if they are not of numerical types.
func Product(top SExpr) (SExpr, error) {
	var err error
	var result SExpr = Integer(1)
	cur := top
	for {
		if IsNull(cur) {
			break
		}
		result, err = Multiply(Car(cur), result)
		if err != nil {
			return nil, err
		}
		cur = Cdr(cur)
	}
	return result, nil
}

// divides the two provided numbers, returns an error if they are not of numerical types.
func Divide(a, b SExpr) (SExpr, error) {
	switch at := a.(type) {
	case Integer:
		switch bt := b.(type) {
		case Integer:
			return Integer(int64(at) / int64(bt)), nil
		case Float:
			return Float(float64(int64(at)) / float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot divide type %T", b)
		}
	case Float:
		switch bt := b.(type) {
		case Integer:
			return Float(float64(at) / float64(int64(bt))), nil
		case Float:
			return Float(float64(at) / float64(bt)), nil
		default:
			return nil, fmt.Errorf("Cannot divide type %T", b)
		}
	default:
		return nil, fmt.Errorf("Cannot divide type %T", a)
	}
}

// divides all the numbers in the provided list. Returns an error if they are not of numerical types.
func Quotient(top SExpr) (SExpr, error) {
	var err error
	if IsNull(top) {
		return nil, fmt.Errorf("division expects at least one parameter")
	}
	var result SExpr = Car(top)
	cur := Cdr(top)
	for {
		if IsNull(cur) {
			break
		}
		result, err = Divide(result, Car(cur))
		if err != nil {
			return nil, err
		}
		cur = Cdr(cur)
	}
	return result, nil
}

