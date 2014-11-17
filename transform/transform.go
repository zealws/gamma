package transform

import (
	"errors"
	"github.com/zfjagann/gamma/sexpr"
)

var (
	Invalid = errors.New("Invalid input SExpr")
)

/*
Type Transform represents a transformation of an SExpr into a different SExpr.

If the Transform is not applicable to the SExpr given, it will return transform.Invalid.
*/
type Transform func(input sexpr.SExpr) (output sexpr.SExpr, err error)

/*
Canonicalize is a Transformation which canonicalizes Closures within an SExpr.
*/
func Canonicalize(input sexpr.SExpr) (sexpr.SExpr, error) {
	return nil, Invalid
}
