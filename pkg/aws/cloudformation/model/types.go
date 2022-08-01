package model

import (
   "fmt"
   "github.com/vektah/gqlparser/v2/ast"
   "strings"
)

var knownTypes = make(map[string]interface{})

func addType(name string) (err error) {
   // Remove not null indicator
   name = strings.TrimSuffix(name, "!")

   if _, found := knownTypes[name]; found {
      return fmt.Errorf("duplicate type: %v", name)
   }
   knownTypes[name] = nil
   return
}

// Helper
func nameFromType(t *ast.Type) string {
   if t.Elem != nil {
      return t.Elem.NamedType
   }
   return t.NamedType
}
