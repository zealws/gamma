package interp

import (
    "fmt"

	. "github.com/zfjagann/gamma/sexpr"
)

func FMap(args SExpr) (SExpr, error) {
	var err error
	if IsNull(top) || IsNull(Cdr(top)) {
		return nil, fmt.Errorf("map expects at exactly two parameters")
	}
    var f SExpr = Car(args)
	var result SExpr = Null
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
