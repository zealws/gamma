package interp

import (
	"fmt"

	. "github.com/zfjagann/gamma/sexpr"
)

func randLength(randList SExpr) int {
	l := 0
	for {
		if IsNull(randList) {
			return l
		}
		l += 1
		randList = Cdr(randList)
	}
}

func checkLen(size int, rator, randList SExpr) error {
	actual := randLength(randList)
	if actual != size {
		return fmt.Errorf("%v expects %d arguments but was given %d", rator, size, actual)
	}
	return nil
}
