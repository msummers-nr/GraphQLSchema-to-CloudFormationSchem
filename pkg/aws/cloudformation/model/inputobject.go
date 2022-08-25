package model

import (
	"github.com/vektah/gqlparser/v2/ast"
)

func NewInputObject(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
	// if added reference or is an array, don't change Property properties
	if property.Ref != "" || property.IsArray {
		return
	}
	property.Type = "object"
	f := new(bool)
	*f = false
	property.AdditionalProperties = f
	return
}
