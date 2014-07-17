package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    filename := os.Args[1]
    f, _ := os.Open(filename)
    expr, _ := ParseSExpr(bufio.NewReader(f))
    fmt.Println(expr.FormatString())
}
