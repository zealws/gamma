go-gamma
========

Golang + Scheme

Submodules
----------

Gamma is divided into separate submodules.

- `gamma` (the root of the project) is the gamma executable binary, and includes the main function for the gamma command line interpreter.
- `gamma/sexpr` includes the type hierarchy for the gamma representation of [s-expressions](http://en.wikipedia.org/wiki/S-expression).
- `gamma/parse` includes the gamma s-expression parsing library.
- `gamma/interp` includes the interpreter implementation.
- `gamma/transform` includes an s-expression manipulation library.

Future Work
-----------

- Compiler?
- Support for concurrency
- Automatically memoized functions
