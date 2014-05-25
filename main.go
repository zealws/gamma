package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    filename := os.Args[1]
    f, _ := os.Open(filename)
    s, _ := ioutil.ReadAll(f)
    expr, _ := ParseSExprString(string(s))
    fmt.Println(expr.FormatString())
}
