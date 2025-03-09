package parser

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"unicode"

	"gotoproto/pkg/models"
)

const (
	tpl = `package main

func main() {
	%s
}`
)

func Parse(s string) (result []models.StructInfo, err error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}

	s = fmt.Sprintf(tpl, s)

	f, err := parser.ParseFile(token.NewFileSet(), "editor.go", s, parser.SpuriousErrors)
	if err != nil {
		return result, fmt.Errorf("failed to parse struct: %w", err)
	}
	result, err = processFile(f)
	if err != nil {
		return result, fmt.Errorf("failed to process file: %w", err)
	}

	return result, nil
}

func processFile(node *ast.File) (allStructs []models.StructInfo, err error) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			structType, ok := x.Type.(*ast.StructType)
			if !ok {
				return false
			}

			structName := x.Name.Name

			allStructs = append(allStructs, models.StructInfo{
				Name:   structName,
				Fields: getFlatFields(structType),
			})
		}
		return true
	})

	if len(allStructs) == 0 {
		return nil, fmt.Errorf("no structs found")
	}

	return allStructs, nil
}

func getTypeSpec(field *ast.Field) (*ast.TypeSpec, bool) {
	ident, isIdentifier := field.Type.(*ast.Ident)
	if !isIdentifier || ident.Obj == nil {
		return nil, false
	}

	typeSpec, isTypeSpec := ident.Obj.Decl.(*ast.TypeSpec)
	if !isTypeSpec {
		return nil, false
	}

	return typeSpec, true
}

func getSimpleType(s string) string {
	switch s {
	case "int", "int64":
		return "int64"
	case "int32":
		return "int32"
	case "float64":
		return "double"
	case "float32":
		return "float"
	case "uint", "uint64":
		return "uint64"
	case "uint32":
		return "uint32"
	case "[]byte", "any":
		return "bytes"
	}

	return s
}

func getType(field any) models.Type {
	switch t := field.(type) {
	case *ast.Ident:
		return models.Type{Name: getSimpleType(t.Name)}
	case *ast.StarExpr:
		if ident, isIdent := t.X.(*ast.Ident); isIdent {
			return models.Type{Name: getSimpleType(ident.Name)}
		}

		return models.Type{Name: "bytes"}
	case *ast.ArrayType:
		if ident, ok := t.Elt.(*ast.Ident); ok {
			simpleType := getSimpleType(ident.String())
			if simpleType == "bytes" {
				return models.Type{Name: "bytes"}
			}

			return models.Type{Name: simpleType + "[]"}
		}

		return models.Type{Name: "bytes"}
	case *ast.SelectorExpr:
		typeName := fmt.Sprintf("%s.%s", t.X, t.Sel)
		switch typeName {
		case "time.Time":
			return models.Type{Name: "string"}
		default:
			return models.Type{CustomType: &typeName}
		}
	case *ast.MapType:
		return models.Type{
			MapType: &models.MapType{
				KeyType:   getType(t.Key).Name,
				ValueType: getType(t.Value).Name,
			},
		}
	case *ast.InterfaceType:
		return models.Type{Name: "bytes"}
	case *ast.StructType:
		return models.Type{
			Fields: getFlatFields(t),
		}
	default:
		return models.Type{Name: "bytes"}
	}
}

func getFlatFields(obj *ast.StructType) (fields []models.Field) {
	for _, field := range obj.Fields.List {
		if len(field.Names) == 0 { // Anonymous field (embedded struct)
			typeSpec, success := getTypeSpec(field)
			if !success {
				continue
			}
			embeddedStruct, ok := typeSpec.Type.(*ast.StructType)
			if ok {
				embeddedFields := getFlatFields(embeddedStruct)
				fields = append(fields, embeddedFields...)
			}

			continue
		}

		// skip unexported fields
		if unicode.IsLower(rune(field.Names[0].Name[0])) {
			continue
		}

		structField := models.Field{
			Name: field.Names[0].Name,
			Type: getType(field.Type),
		}

		fields = append(fields, structField)
	}

	return fields
}
