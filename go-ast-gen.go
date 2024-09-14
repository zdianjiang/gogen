package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Model struct {
	ProjectName  string
	ModelName    string
	IdType       string
	Fields       []Field
	ModelPackage string
	Relations    []Relation
}

type Field struct {
	Name string
	Type string
}

type Relation struct {
	Name string
	Type string
}

func main() {
	// Define command line arguments
	modelFile := flag.String("model", "", "Path to the model file")
	templateFile := flag.String("template", "", "Path to the template file")
	outputFile := flag.String("output", "", "Path to the output file")
	projectName := flag.String("project", "", "Name of the project")
	flag.Parse()

	if *modelFile == "" || *templateFile == "" || *outputFile == "" {
		fmt.Println("Usage: go run main.go -model <model_file> -template <template_file> -output <output_file>")
		return
	}

	// Determine the model package import path
	modelPackage := determineModelPackage(*modelFile)

	// Parse the model file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *modelFile, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing model file: %v\n", err)
		return
	}

	// Extract model information
	var model Model
	model.ModelPackage = modelPackage
	model.ProjectName = *projectName
	fmt.Println("project:", model.ProjectName)
	fmt.Println("modelPackage:", model.ModelPackage)
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						model.ModelName = typeSpec.Name.Name
						for _, field := range structType.Fields.List {
							for _, name := range field.Names {
								fieldType := fieldTypeToString(field.Type)
								model.Fields = append(model.Fields, Field{
									Name: name.Name,
									Type: fieldType,
								})
								if name.Name == "ID" {
									model.IdType = fieldType
								}
								// Check for relations
								if strings.HasPrefix(fieldType, "[]") || strings.HasPrefix(fieldType, "*") {
									model.Relations = append(model.Relations, Relation{
										Name: name.Name,
										Type: fieldType,
									})
								}
							}
						}
						return false
					}
				}
			}
		}
		return true
	})

	// Parse the template file
	tmpl, err := template.ParseFiles(*templateFile)
	if err != nil {
		fmt.Printf("Error parsing template file: %v\n", err)
		return
	}

	// Execute the template with the model data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, model)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Format the generated code
	formattedSrc, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("Error formatting source: %v\n", err)
		fmt.Println("Generated source:")
		fmt.Println(buf.String())
		return
	}

	// Create the output file
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the formatted code to the output file
	_, err = file.Write(formattedSrc)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}
}

func determineModelPackage(modelFile string) string {
	absPath, err := filepath.Abs(modelFile)
	if err != nil {
		panic(err)
	}
	// Assuming the project is within the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	relPath, err := filepath.Rel(wd, absPath)
	if err != nil {
		panic(err)
	}
	// Remove the filename to get the package path
	packagePath := filepath.Dir(relPath)
	return strings.ReplaceAll(packagePath, string(os.PathSeparator), "/")
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
