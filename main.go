package main

import (
	"github.com/zfjagann/go-gamma/interp"
	"github.com/zfjagann/go-gamma/parse"

	"fmt"
	"os"
)

func main() {
	input := os.Stdin
	parser := parse.NewParser(input)
	eval := interp.NewInterpreter(interp.DefaultEnvironment)
	for {
		fmt.Print("scheme00> ")
		input, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
			fmt.Println()
		}
		output, err := eval.Evaluate(input)
		if err != nil {
			if err == interp.Exit {
				return
			}
			fmt.Println(err)
		} else {
			fmt.Printf("%v\n", output)
		}
	}
	os.Exit(0)
}
