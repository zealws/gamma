package main

import (
	"flag"
	"fmt"
	"github.com/zfjagann/gamma/interp"
	"github.com/zfjagann/gamma/parse"
	"github.com/zfjagann/gamma/sexpr"
	"io"
	"os"
)

func main() {
	fname := flag.String("f", "-", "specify a file to run")
	flag.Parse()

	if *fname == "-" {
		os.Exit(repl(true, os.Stdin))
	} else {
		input, err := os.Open(*fname)
		if err != nil {
			fmt.Println(err)
			os.Exit(255)
		}
		os.Exit(repl(false, input))
	}
}

func repl(interactive bool, input io.Reader) int {
	parser := parse.NewParser(input)
	eval := interp.NewInterpreter(interp.DefaultEnvironment)
	for {
		if interactive {
			fmt.Print("scheme00> ")
		}
		input, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				return 0
			}
			fmt.Println(err)
			if !interactive {
				return 1
			}
			fmt.Println()
			continue
		}
		output, err := eval.Evaluate(input)
		if err != nil {
			if err == interp.Exit {
				return 0
			}
			fmt.Println(err)
			if !interactive {
				return 2
			}
		} else if output != sexpr.Null {
			fmt.Printf("%v\n", output)
		}
	}
}

func run(fname string) {
	input, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	parser := parse.NewParser(input)
	eval := interp.NewInterpreter(interp.DefaultEnvironment)
	for {
		input, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			fmt.Println(err)
			os.Exit(2)
		}
		output, err := eval.Evaluate(input)
		if err != nil {
			if err == interp.Exit {
				return
			}
			fmt.Println(err)
			os.Exit(3)
		} else {
			fmt.Printf("%v\n", output)
		}
	}
}
