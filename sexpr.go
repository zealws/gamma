package main

import (
    "bufio"
    "fmt"
    "io"
    "strconv"
    "strings"
)

const (
    whitespace = " \n\r\t"
)

var (
    Null = nullExpr{}
)

type SExpr interface {
    FormatString() string
    Equals(SExpr) bool
}

/**
*** Null Type
**/

// Singleton - for simplicity of comparison.
// Use Null variable defined above.
type nullExpr struct{}

func (x nullExpr) FormatString() string {
    return "_"
}

func (x nullExpr) Equals(other SExpr) bool {
    return other == Null
}

/**
*** PairExpr
**/

type PairExpr struct {
    car SExpr
    cdr SExpr
}

func (x PairExpr) Car() SExpr {
    return x.car
}

func (x PairExpr) Cdr() SExpr {
    return x.cdr
}

func (x PairExpr) FormatString() string {
    cars, cdrs := "!", "!"
    if x.car != nil {
        cars = x.car.FormatString()
    }
    if x.cdr != nil {
        cdrs = x.cdr.FormatString()
    }
    return fmt.Sprintf("(%s . %s)", cars, cdrs)
}

func (x PairExpr) Equals(other SExpr) bool {
    switch o := other.(type) {
    case *PairExpr:
        return x.car.Equals(o.car) && x.cdr.Equals(o.cdr)
    case PairExpr:
        return x.car.Equals(o.car) && x.cdr.Equals(o.cdr)
    default:
        return false
    }
}

/**
*** IntExpr
**/

type IntExpr struct {
    value int64
}

func (x IntExpr) Int() int64 {
    return x.value
}

func (x IntExpr) FormatString() string {
    return fmt.Sprintf("%d", x.Int())
}

func (x IntExpr) Equals(other SExpr) bool {
    switch o := other.(type) {
    case *IntExpr:
        return x.value == o.value
    case IntExpr:
        return x.value == o.value
    default:
        return false
    }
}

/**
*** SymbolExpr
**/

type SymbolExpr struct {
    name string
}

func (x SymbolExpr) Name() string {
    return x.name
}

func (x SymbolExpr) FormatString() string {
    return x.name
}

func (x SymbolExpr) Equals(other SExpr) bool {
    switch o := other.(type) {
    case *SymbolExpr:
        return x.name == o.name
    case SymbolExpr:
        return x.name == o.name
    default:
        return false
    }
}

/**
*** StringExpr
**/

type StringExpr struct {
    value string
}

func (x StringExpr) String() string {
    return x.value
}

func (x StringExpr) FormatString() string {
    return fmt.Sprintf(`"%s"`, x.String())
}

func (x StringExpr) Equals(other SExpr) bool {
    switch o := other.(type) {
    case *StringExpr:
        return x.value == o.value
    case StringExpr:
        return x.value == o.value
    default:
        return false
    }
}

/**
*** Parsing
**/

func ParseSExprString(in string) (SExpr, error) {
    return ParseSExpr(bufio.NewReader(strings.NewReader(in)))
}

func ParseSExpr(in *bufio.Reader) (SExpr, error) {
    err := discardWhitespace(in)
    if err != nil {
        return nil, err
    }
    next, err := in.Peek(1)
    if err != nil {
        return nil, err
    }
    if strings.Contains("0123456789", string(next)) {
        return parseIntExpr(in)
    } else if string(next) == "\"" {
        return parseStringExpr(in)
    } else if string(next) == "(" {
        return parsePairExpr(in)
    } else if string(next) == "[" {
        return parseListExpr(in)
    } else {
        return parseSymbolExpr(in)
    }
}

func parseStringExpr(in *bufio.Reader) (SExpr, error) {
    // First char MUST be "
    in.ReadRune()
    str := ""
    for {
        c, err := in.Peek(1)
        if err == io.EOF {
            return nil, NewError("EOF while reading string constant")
        } else if err != nil {
            return nil, err
        } else if string(c) == "\"" {
            in.ReadRune()
            break
        } else {
            r, _, err := in.ReadRune()
            if err != nil {
                return nil, err
            }
            str += string(r)
        }
    }
    return &StringExpr{str}, nil
}

func parseIntExpr(in *bufio.Reader) (SExpr, error) {
    str := ""
    for {
        c, err := in.Peek(1)
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        } else if strings.Contains(whitespace, string(c)) {
            break
        } else if string(c) == ")" {
            break
        } else if string(c) == "]" {
            break
        } else {
            r, _, err := in.ReadRune()
            if err != nil {
                return nil, err
            }
            str += string(r)
        }
    }
    s, err := strconv.ParseInt(str, 10, 64)
    if err != nil {
        return nil, err
    }
    return &IntExpr{s}, nil
}

func parseSymbolExpr(in *bufio.Reader) (SExpr, error) {
    name := ""
    for {
        c, err := in.Peek(1)
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        } else if strings.Contains(whitespace, string(c)) {
            break
        } else if string(c) == ")" {
            break
        } else if string(c) == "]" {
            break
        } else {
            r, _, err := in.ReadRune()
            if err != nil {
                return nil, err
            }
            name += string(r)
        }
    }
    if name == "_" {
        return Null, nil
    }
    return &SymbolExpr{name}, nil
}

func parseListExpr(in *bufio.Reader) (SExpr, error) {
    err := readRune('[', in)
    if err != nil {
        return nil, err
    }
    var expr, curr *PairExpr

    for {
        c, err := in.Peek(1)
        if err == io.EOF {
            return nil, NewError("EOF while reading array")
        } else if err != nil {
            return nil, err
        } else if string(c) == "]" {
            break
        } else {
            if curr == nil {
                expr = &PairExpr{}
                expr.car, err = ParseSExpr(in)
                curr = expr
            }
            next := &PairExpr{}
            curr.cdr = next
            curr = next
            curr.car, err = ParseSExpr(in)
            if err != nil {
                return nil, err
            }
        }
    }
    if curr != nil {
        curr.cdr = Null
    }

    err = readRune(']', in)
    if err != nil {
        return nil, err
    }
    if expr != nil {
        return expr, nil
    } else {
        return Null, nil
    }
}

func parsePairExpr(in *bufio.Reader) (SExpr, error) {
    err := readRune('(', in)
    if err != nil {
        return nil, err
    }
    car, err := ParseSExpr(in)
    if err != nil {
        return nil, err
    }
    cdr, err := parseCdr(in)
    if err != nil {
        return nil, err
    }
    err = readRune(')', in)
    if err != nil {
        return nil, err
    }
    return &PairExpr{car, cdr}, nil
}

func parseCdr(in *bufio.Reader) (SExpr, error) {
    // ". (b . c))"
    err := discardWhitespace(in)
    if err != nil {
        return nil, err
    }
    err = readRune('.', in)
    if err != nil {
        return nil, err
    }
    return ParseSExpr(in)
}

func discardWhitespace(in *bufio.Reader) (err error) {
    x, _, err := in.ReadRune()
    if err != nil {
        return
    }
    for strings.Contains(whitespace, string(x)) {
        x, _, err = in.ReadRune()
        if err != nil {
            return
        }
    }
    in.UnreadRune()
    return
}

func readRune(r rune, in *bufio.Reader) error {
    c, _, err := in.ReadRune()
    if err != nil {
        return err
    }
    if rune(c) != r {
        return NewError("Found", string(c), "instead of", string(r))
    }
    return nil
}
