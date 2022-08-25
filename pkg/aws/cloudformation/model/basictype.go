package model

import (
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

/*
JSON Schema types:
string.
number. //float
integer. //int
object.
array.
boolean.
null.
*/

// NewBasicType Represent basic, built-in (in the language sense) types
func NewBasicType(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
	// if is an array, don't change Property properties
	if property.IsArray {
		return
	}

	//translate GraphQL to JSON basic types
	switch typeDef.NamedType {
	case "Float":
		property.Type = "number"
	case "Int":
		property.Type = "integer"
	case "ID":
		property.Type = "string"
	default:
		property.Type = strings.ToLower(typeDef.NamedType)
	}

	return
}

func NewBasicTypeDefinition(name string, directives ast.DirectiveList) *ast.Definition {
	return &ast.Definition{
		Description: "",
		Name:        name,
		Directives:  directives,
		Interfaces:  nil,
		Fields:      nil,
		EnumValues:  nil,
		Position:    nil,
		BuiltIn:     false,
	}
}
