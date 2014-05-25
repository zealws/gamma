package main

import (
    "testing"
)

func checkValue(t *testing.T, x interface{}, err error) {
    if err != nil {
        t.Fatal(err)
    }
    if x == nil {
        t.Fatal("Nil returned")
    }
}

/**
*** IntExpr Assertions
**/

func assertIsInteger(t *testing.T, x SExpr) *IntExpr {
    if x == nil {
        t.Fatal("Expected Int but was null.")
    }
    switch xpr := x.(type) {
    case *IntExpr:
        return xpr
    default:
        t.Fatalf("Non-integer %T returned: %s", xpr, xpr.FormatString())
        return nil
    }
}

func assertIntegerValue(t *testing.T, i *IntExpr, value int64) {
    if i.Int() != value {
        t.Fatalf("Integer does not match expected %d != %d", i.Int(), value)
    }
}

/**
*** SymbolExpr Assertions
**/

func assertIsSymbol(t *testing.T, x SExpr) *SymbolExpr {
    if x == nil {
        t.Fatal("Expected Symbol but was null.")
    }
    switch xpr := x.(type) {
    case *SymbolExpr:
        return xpr
    default:
        t.Fatalf("Non-symbol %T returned: %s", xpr, xpr.FormatString())
        return nil
    }
}

func assertSymbolName(t *testing.T, x *SymbolExpr, name string) {
    if x.Name() != name {
        t.Fatalf("Symbol name does not match expected %s != %s", x.Name(), name)
    }
}

/**
*** StringExpr Assertions
**/

func assertIsString(t *testing.T, x SExpr) *StringExpr {
    if x == nil {
        t.Fatal("Expected String but was null.")
    }
    switch xpr := x.(type) {
    case *StringExpr:
        return xpr
    default:
        t.Fatalf("Non-symbol %T returned: %s", xpr, xpr.FormatString())
        return nil
    }
}

func assertStringValue(t *testing.T, xpr *StringExpr, value string) {
    if xpr.String() != value {
        t.Fatalf(`String does not match expected %s != "asdf"`, xpr.String())
    }
}

/**
*** PairExpr Assertions
**/

func assertIsPair(t *testing.T, x SExpr) *PairExpr {
    if x == nil {
        t.Fatal("Expected Pair but was null.")
    }
    switch xpr := x.(type) {
    case *PairExpr:
        return xpr
    default:
        t.Fatalf("Non-pair %T returned: %s", xpr, xpr.FormatString())
        return nil
    }
}

/**
*** Null Assertions
**/

func assertIsNull(t *testing.T, x SExpr) {
    if x != Null {
        t.Fatalf("Non-null %T returned: %s", x, x.FormatString())
    }
}

/**
*** Tests
**/

func TestParseInteger(t *testing.T) {
    x, err := ParseSExprString("1234567")
    checkValue(t, x, err)
    i := assertIsInteger(t, x)
    assertIntegerValue(t, i, 1234567)
}

func TestParseString(t *testing.T) {
    x, err := ParseSExprString(`"asdf"`)
    checkValue(t, x, err)
    str := assertIsString(t, x)
    assertStringValue(t, str, "asdf")
}

func TestParseSymbol(t *testing.T) {
    x, err := ParseSExprString("asdf")
    checkValue(t, x, err)
    sym := assertIsSymbol(t, x)
    assertSymbolName(t, sym, "asdf")
}

func TestStripsWhitespaceBeforeString(t *testing.T) {
    x, err := ParseSExprString(`         "    asdf    "`)
    checkValue(t, x, err)
    str := assertIsString(t, x)
    assertStringValue(t, str, "    asdf    ")
}

func TestParsesPairs(t *testing.T) {
    x, err := ParseSExprString(`(1357 . 2468)`)
    checkValue(t, x, err)
    pair := assertIsPair(t, x)
    car, cdr := pair.Car(), pair.Cdr()
    assertIntegerValue(t, assertIsInteger(t, car), 1357)
    assertIntegerValue(t, assertIsInteger(t, cdr), 2468)

    x, err = ParseSExprString(`(abcd . (_ . efgh))`)
    checkValue(t, x, err)
    pair = assertIsPair(t, x)
    car, cdr = pair.Car(), pair.Cdr()
    cdrp := assertIsPair(t, cdr)
    assertSymbolName(t, assertIsSymbol(t, car), "abcd")
    cadr, cddr := cdrp.Car(), cdrp.Cdr()
    assertIsNull(t, cadr)
    assertSymbolName(t, assertIsSymbol(t, cddr), "efgh")

    x, err = ParseSExprString(`((1234 . "asdf") . (_ . xyz))`)
    checkValue(t, x, err)
    pair = assertIsPair(t, x)
    car, cdr = pair.Car(), pair.Cdr()
    carp, cdrp := assertIsPair(t, car), assertIsPair(t, cdr)
    caar, cdar := carp.Car(), carp.Cdr()
    cadr, cddr = cdrp.Car(), cdrp.Cdr()
    assertIntegerValue(t, assertIsInteger(t, caar), 1234)
    assertStringValue(t, assertIsString(t, cdar), "asdf")
    assertIsNull(t, cadr)
    assertSymbolName(t, assertIsSymbol(t, cddr), "xyz")
}

func TestParsesLists(t *testing.T) {
    x, err := ParseSExprString(`[1234 5678]`)
    checkValue(t, x, err)
    o, err := ParseSExprString("(1234 . (5678 . _))")
    checkValue(t, o, err)
    if !x.Equals(o) {
        t.Fatal("List did not match expected contents: [1234 5678]")
    }

    x, err = ParseSExprString(`[]`)
    checkValue(t, x, err)
    o, err = ParseSExprString("_")
    checkValue(t, o, err)
    if !x.Equals(o) {
        t.Fatal("List did not match expected contents: []")
    }
}
