package model

import (
	"github.com/vektah/gqlparser/v2/ast"
	"reflect"
)

func NewEnum(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
	// if added reference or is an array, don't change Property properties
	if property.Ref != "" || property.IsArray { //if already written in definition
		return
	}
	//evaluates to string
	property.Type = reflect.TypeOf(definition.EnumValues[0].Name).String()

	//list enum values under "enum" property
	enumValues := make([]string, 0)
	for _, value := range definition.EnumValues {
		enumValues = append(enumValues, value.Name)
	}
	property.Enum = enumValues
	return
}
