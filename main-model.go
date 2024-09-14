package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func main() {
	src := `package model; type User struct { ID int; Name string; Age int; }`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Find the User struct and generate create function
	var createFuncCode string
	ast.Inspect(f, func(n ast.Node) bool {
		// Find type declarations
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						if typeSpec.Name.Name == "User" {
							createFuncCode = generateCreateFunc(typeSpec.Name.Name, structType)
							return false
						}
					}
				}
			}
		}
		return true
	})

	if createFuncCode != "" {
		fmt.Println(createFuncCode)
	} else {
		log.Println("User struct not found in src")
	}
}

func generateCreateFunc(structName string, structType *ast.StructType) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("func create%s(", structName))

	// Generate function parameters
	params := []string{}
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			param := fmt.Sprintf("%s %s", name.Name, fieldTypeToString(field.Type))
			params = append(params, param)
		}
	}
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(fmt.Sprintf(") *%s {\n", structName))

	// Generate function body
	buf.WriteString(fmt.Sprintf("\treturn &%s{\n", structName))
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			buf.WriteString(fmt.Sprintf("\t\t%s: %s,\n", name.Name, name.Name))
		}
	}
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	return buf.String()
}

func fieldTypeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", fieldTypeToString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", fieldTypeToString(t.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", fieldTypeToString(t.Elt))
	default:
		return ""
	}
}
