package main

import (
	"fmt"
	"os"

	"github.com/dpinela/regexc/ast"
	"github.com/k0kubun/pp"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Please specify a regexp.")
		return
	}
	re := os.Args[1]
	pp.Print(ast.Parse(re))
}
