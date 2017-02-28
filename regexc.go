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
	tree, err := ast.Parse(re)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	pp.Print(tree)
}
