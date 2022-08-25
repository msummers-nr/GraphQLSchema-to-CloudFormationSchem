package model

import (
	"github.com/vektah/gqlparser/v2/ast"
)

// GraphQL Scalars are leaf nodes that are implementation dependent, not "scalar" in the programming language sense.
//
// See: https://graphql.org/learn/schema/#scalar-types

func NewScalar(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
	// if added reference or is an array, don't change Property properties
	if property.IsArray || property.Ref != "" {
		return
	}
	property.Type = "string"

	return
}
