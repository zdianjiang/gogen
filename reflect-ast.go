package main

import (
	"go/ast"
	"go/printer"
	"go/token"
	"main/model"
	"os"
	"path"
	"reflect"
)

func parse2ast(t reflect.Type) ast.Expr {
	if t.Kind() != reflect.Struct {
		panic("not a struct type")
	}

	structType := &ast.StructType{Struct: token.Pos(1), Fields: &ast.FieldList{}}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fieldIdent := &ast.Ident{Name: field.Name}
		var fieldType ast.Expr

		switch field.Type.Kind() {
		case reflect.Struct:
			fieldType = parse2ast(field.Type)
		case reflect.Slice:
			elemType := field.Type.Elem()
			if elemType.Kind() == reflect.Struct {
				fieldType = &ast.ArrayType{
					Elt: &ast.SelectorExpr{
						X:   &ast.Ident{Name: path.Base(elemType.PkgPath())},
						Sel: &ast.Ident{Name: elemType.Name()},
					},
				}
			} else {
				fieldType = &ast.ArrayType{
					Elt: &ast.Ident{Name: elemType.Name()},
				}
			}
		default:
			if field.Type.PkgPath() != "" {
				fieldType = &ast.SelectorExpr{
					X:   &ast.Ident{Name: path.Base(field.Type.PkgPath())},
					Sel: &ast.Ident{Name: field.Type.Name()},
				}
			} else {
				fieldType = &ast.Ident{Name: field.Type.Name()}
			}
		}

		fieldAST := &ast.Field{
			Names: []*ast.Ident{fieldIdent},
			Type:  fieldType,
		}

		structType.Fields.List = append(structType.Fields.List, fieldAST)
	}

	return structType
}

func main() {
	userType := reflect.TypeOf(model.User{})
	userAST := parse2ast(userType)

	// // 获取包的导入路径
	// pkgPath := userType.PkgPath()
	// // 打印包的导入路径
	// fmt.Println("Package Path:", pkgPath)
	// // 使用 Name() 方法获取结构体的名称
	// structName := userType.Name()

	// 打印生成的 AST
	ast.Print(nil, userAST)

	fset := token.NewFileSet()
	printer.Fprint(os.Stdout, fset, userAST)

	// 创建一个 *ast.File 节点，因为 printer.Fprint 需要一个 *ast.File 作为参数
	// file := &ast.File{
	// 	Name: &ast.Ident{Name: pkgPath},
	// 	Decls: []ast.Decl{&ast.GenDecl{
	// 		Tok: token.TYPE,
	// 		Specs: []ast.Spec{&ast.TypeSpec{
	// 			Name: &ast.Ident{Name: structName},
	// 			Type: userAST,
	// 		}},
	// 	}},
	// }

	// // 使用 go/format 包来格式化 AST 为 Go 源代码
	// var buf bytes.Buffer
	// if err := format.Node(&buf, token.NewFileSet(), file); err != nil {
	// 	panic(err)
	// }

	// // 输出格式化后的 Go 源代码
	// src := buf.Bytes()
	// println(string(src))
}
