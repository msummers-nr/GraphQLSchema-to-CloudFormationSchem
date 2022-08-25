package model

import (
	"github.com/vektah/gqlparser/v2/ast"
)

func NewUnion(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
	// if added reference or is an array, don't change Property properties
	if property.Ref != "" || property.IsArray {
		return
	}
	unionItem := Item{}

	// for each type in the union type, create Item with $ref property
	// should be anyOf or oneOf?
	for _, unionType := range definition.Types {
		unionTypes := Item{}
		unionTypes.Ref = "#/definitions/" + unionType
		unionItem.AnyOf = append(unionItem.AnyOf, &unionTypes)
	}

	property.Type = "object"
	property.AnyOf = unionItem.AnyOf
	f := new(bool)
	*f = false
	property.AdditionalProperties = f
	return
}
