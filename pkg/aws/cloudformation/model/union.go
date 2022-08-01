package model

import (
   "github.com/vektah/gqlparser/v2/ast"
)

// TODO
func NewUnion(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
   property.Type = "// FIXME Union"
   return
}
