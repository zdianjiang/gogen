package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	src := `
package hello

import "fmt"

func greet() {
	fmt.Println("Hello World!")
}
`
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}
