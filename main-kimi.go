package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// 模型定义的源代码
	src := `package model; type User struct { ID int; Name string; Age int; }`
	// 设置文件集
	fset := token.NewFileSet()
	// 解析源代码
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// 遍历AST
	ast.Inspect(file, func(n ast.Node) bool {
		// 查找类型声明
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				// 查找类型定义
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					// 这里我们假设模型名称和类型名称相同
					modelName := typeSpec.Name.Name
					fmt.Printf("Creating GORM code for model: %s\n", modelName)

					// 生成创建表的代码
					fmt.Printf("gorm.Model{} // GORM embedded model for common fields\n")
					// 遍历结构体字段
					structType, ok := typeSpec.Type.(*ast.StructType)
					if ok {
						for _, field := range structType.Fields.List {
							// 这里我们简单地处理每个字段
							for _, name := range field.Names {
								fmt.Printf("%s %s `gorm:\"type:%s;\"`\n",
									name.Name, field.Type, field.Type.(*ast.Ident).Name)
							}
						}
					}
				}
			}
		}
		return true
	})
}
