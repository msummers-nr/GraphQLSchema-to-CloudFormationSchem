package model

import (
   "github.com/vektah/gqlparser/v2/ast"
   "strings"
)

/*
JSON Schema types:
string.
number.
integer.
object.
array.
boolean.
null.
*/

// Represent basic, built-in (in the language sense) types
func NewBasicType(property *Property, definition *ast.Definition, typeDef *ast.Type) (err error) {
   // log.Printf("NewBasicType: ast.Type: %+V", typeDef)
   property.Type = strings.ToLower(typeDef.NamedType)
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

// TODO translate GraphQL scalars to JSON scalars
