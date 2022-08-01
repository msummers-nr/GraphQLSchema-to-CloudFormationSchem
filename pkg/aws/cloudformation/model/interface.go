package model

import (
   "github.com/vektah/gqlparser/v2/ast"
)

// TODO
func NewInterface(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
   property.Type = "// FIXME Interface"
   return
}
