package model

import (
   "fmt"
   log "github.com/sirupsen/logrus"
   "github.com/vektah/gqlparser/v2/ast"
)

type Document struct {
   TypeName             string               `json:"typeName"`
   Description          string               `json:"description"`
   SourceUrl            string               `json:"sourceUrl"`
   Definitions          map[string]*Property `json:"definitions"`
   Properties           map[string]*Property `json:"properties"`
   AdditionalProperties bool                 `json:"additionalProperties"`
   Required             []string             `json:"required"`
   ReadOnlyProperties   []string             `json:"readOnlyProperties"`
   PrimaryIdentifier    []string             `json:"primaryIdentifier"`
   Handlers             map[string]*Handler  `json:"handlers"`
}

type Handler struct {
}

func NewDocument() (document *Document) {
   document = &Document{
      TypeName:             "",
      Description:          "",
      SourceUrl:            "",
      Definitions:          make(map[string]*Property),
      Properties:           make(map[string]*Property),
      AdditionalProperties: false,
      Required:             make([]string, 0),
      ReadOnlyProperties:   make([]string, 0),
      PrimaryIdentifier:    make([]string, 0),
      Handlers:             make(map[string]*Handler),
   }
   return
}

func (d *Document) AddDefinition(astType *ast.Type, property *Property) (err error) {
   if err = addType(astType.NamedType); err != nil {
      return err
   }

   if property == nil {
      return fmt.Errorf("cannot add nil property to document")
   }
   d.Definitions[astType.NamedType] = property
   return
}

func (d *Document) AddProperty(argDef *ast.ArgumentDefinition, property *Property) (err error) {
   if property == nil {
      return fmt.Errorf("cannot add nil property to document")
   }
   d.Properties[argDef.Name] = property
   if property.IsRequired {
      d.Required = append(d.Required, argDef.Name)
   }
   return
}

// Recursively travel down the astType
func (d *Document) SplunkTypeDefinitions(astType *ast.Type, gqlSchema *ast.SchemaDocument) {
   log.Printf("SplunkDefintions: astType.NamedType: %v", astType.NamedType)
   var def *ast.Definition
   // If it's an array the actual type is Elem
   if astType.Elem != nil {
      astType = astType.Elem
   }
   def = gqlSchema.Definitions.ForName(astType.NamedType)
   if def == nil {
      return
   }
   property, err := NewDefinitionProperty(def, astType)
   if err != nil {
      log.Errorf("SplunkTypeDefinitions: error creating property: %v", err)
      return
   }

   for _, field := range def.Fields {
      fieldTypeName := nameFromType(field.Type)
      fieldTypeDef := gqlSchema.Definitions.ForName(fieldTypeName)
      if fieldTypeDef == nil {
         fieldTypeDef = NewBasicTypeDefinition(field.Name, field.Directives)
      }
      // FIXME Close- this gets the type correct but misses the array if present
      // THis is close, but now there are no arrays
      fieldType := field.Type
      if fieldType.Elem != nil {
         fieldType = fieldType.Elem
      }
      fieldTypeProperty, err := NewProperty(fieldTypeDef, fieldType)
      if err != nil {
         log.Warnf("error getting field type as a property: %v", err)
         continue
      }
      if fieldTypeProperty == nil {
         continue
      }
      property.Properties[field.Name] = fieldTypeProperty
      if field.Type.NonNull {
         property.Required = append(property.Required, field.Name)
      }
      d.SplunkTypeDefinitions(field.Type, gqlSchema)
   }
   err = d.AddDefinition(astType, property)
   if err != nil {
      log.Warnf("error adding definition to document: %v", err)
   }
}
