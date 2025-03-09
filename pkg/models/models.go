package models

type MapType struct {
	KeyType   string
	ValueType string
}

type Type struct {
	Name       string
	CustomType *string
	MapType    *MapType
	Fields     []Field
}

type Field struct {
	Name string
	Type Type
}
type StructInfo struct {
	Name   string
	Fields []Field
}
