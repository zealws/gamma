package main

import (
    "fmt"
    "io"
    "runtime"
)

const (
    // Max size of stack traces in bytes.
    traceback_length = 10240
)

type Error struct {
    message string
    stack   string
}

func (e Error) Error() string {
    return e.message
}

func (e Error) Log(out io.Writer) {
    wrt := func(things ...interface{}) {
        out.Write([]byte(fmt.Sprintln(things...)))
    }
    wrt("gamma.Error:", e.Error())
    wrt(e.stack)
}

func NewError(items ...interface{}) Error {
    tb := make([]byte, traceback_length)
    runtime.Stack(tb, false)
    return Error{fmt.Sprint(items), string(tb)}

}
