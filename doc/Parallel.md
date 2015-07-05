# Gamma Parallelization

## `pexec`

The function `pexec` (short for *parallel execute*) allows you to execute gamma code in parallel.

The `pexec` is invoked as follows:

    (pexec THUNK)

`pexec` will execute `THUNK` in a parallel execution thread, allowing the caller of pexec to continue running.

`THUNK` must be a closure expecting exactly 0 arguments.

`pexec` returns a thunk that can be called to retrieve the result of applying `X`, waiting for `X` to complete if necessary.

Example:

    scheme00> (define x (pexec (sleep 1))
    scheme00> (x)
    1436074148047
