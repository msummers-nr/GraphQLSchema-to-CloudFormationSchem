package model

import (
   "github.com/vektah/gqlparser/v2/ast"
   "golang.org/x/text/cases"
   "golang.org/x/text/language"
   "strings"
)

// Helper
func nameFromType(t *ast.Type) string {
   if t.Elem != nil {
      return t.Elem.NamedType
   }
   return t.NamedType
}

// uppercaseTypeName Helper
// capitalize first letter of property, consistent with cloudformation schema general practice
func uppercaseTypeName(s string) string {
   return cases.Title(language.Und, cases.NoLower).String(s)
}

// IsArrayType Helper
// return if property is an array based on graphql syntax
func IsArrayType(s string) bool {
   return strings.HasPrefix(s, "[")
}
