package model

import (
   "fmt"
   log "github.com/sirupsen/logrus"
   "github.com/vektah/gqlparser/v2/ast"
   "strings"
)

type Document struct {
   TypeName             string                 `json:"typeName"`
   Description          string                 `json:"description"`
   SourceUrl            string                 `json:"sourceUrl,omitempty"` // not required, but must have pattern ^https://[0-9a-zA-Z]([-.\\w]*[0-9a-zA-Z])(:[0-9]*)*([?/#].*)?$
   Definitions          map[string]*Property   `json:"definitions"`
   Properties           map[string]*Property   `json:"properties"`
   AdditionalProperties bool                   `json:"additionalProperties"`
   Required             []string               `json:"required"`
   ReadOnlyProperties   []string               `json:"readOnlyProperties"`
   PrimaryIdentifier    []string               `json:"primaryIdentifier"`
   Handlers             map[string]*Handler    `json:"handlers"`
   Tagging              map[string]interface{} `json:"tagging"`
   knownTypes           map[string]interface{} `json:"-"`
}

type Handler struct {
   Permissions []string `json:"permissions"`
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
      Tagging:              make(map[string]interface{}),
      knownTypes:           make(map[string]interface{}),
   }
   document.ReadOnlyProperties = append(document.ReadOnlyProperties, "/properties/Guid")
   document.PrimaryIdentifier = append(document.PrimaryIdentifier, "/properties/Guid")
   handler := &Handler{}
   permissions := []string{"cloudformation:BatchDescribeTypeConfigurations"}
   handler.Permissions = permissions
   document.Handlers["create"] = handler
   document.Handlers["update"] = handler
   document.Handlers["delete"] = handler
   document.Handlers["read"] = handler
   document.Handlers["list"] = handler
   document.Tagging["taggable"] = true
   document.Tagging["tagProperty"] = "#/definitions/EntityTag"

   return
}

// AddDefinition
// add definition property to document
func (d *Document) AddDefinition(astType *ast.Type, property *Property) (err error) {
   // TODO make this smarter so it does a merge rather than overwriting an existing definition
   // testing purposes
   /*if err = addType(astType.Name(), false); err != nil {
   	return err
   }*/
   if property == nil {
      return fmt.Errorf("cannot add nil property to document")
   }
   d.Definitions[astType.NamedType] = property // add to definitions
   return
}

// AddProperty
// add property to document
func (d *Document) AddProperty(argDefName string, property *Property) (err error) {
   if property == nil {
      return fmt.Errorf("cannot add nil property to document")
   }

   existingProperty := d.Properties[uppercaseTypeName(argDefName)]
   if existingProperty == nil {
      d.Properties[uppercaseTypeName(argDefName)] = property
      // FIXME Required is a problem during merge
      if property.IsRequired {
         d.Required = append(d.Required, uppercaseTypeName(argDefName))
      }
      return
   }

   // If the existingProperty IS NOT A REF
   if existingProperty.Ref == "" {
      if existingProperty.Type != property.Type {
         log.Warnf("AddProperty: skipping attempt to merge incompatible property types. Original %+v New: %+v", existingProperty, property)
      } else {
         // Ignore, they're the same
      }
      return
   }

   // The property IS NOT A REF
   if property.Ref == "" {
      log.Warnf("AddProperty: skipping attempt to merge incompatible property types. Original %+v New: %+v", existingProperty, property)
      return
   }

   ref, _ := strings.CutPrefix(existingProperty.Ref, "#/definitions/")
   existingDef := d.Definitions[ref]
   if existingDef == nil {
      log.Warnf("AddProperty: reference not found: %s", existingProperty.Ref)
      return
   }

   ref, _ = strings.CutPrefix(property.Ref, "#/definitions/")
   propertyDef := d.Definitions[ref]
   if propertyDef == nil {
      log.Warnf("AddProperty: reference not found: %s", property.Ref)
      return
   }

   d.mergeProperties(propertyDef, existingDef)
   return
}

func (d *Document) getActualProperty(p *Property) *Property {
   // Not a reference
   if p.Ref == "" {
      return p
   }
   ref, _ := strings.CutPrefix(p.Ref, "#/definitions/")
   refProperty, ok := d.Definitions[ref]
   if !ok {
      log.Warnf("getActualProperty: reference not found: %s", p.Ref)
      return p
   }
   return refProperty
}

func (d *Document) mergeProperties(source *Property, destination *Property) {
   // FIXME don't merge a Property whose name is Properties
   if source.Name == "Properties" {
      return
   }
   // source := d.getActualProperty(src)
   // destination := d.getActualProperty(dst)
   for key, property := range source.Properties {
      if key == "Properties" {
         continue
      }
      if destination.Properties[key] == nil {
         destination.Properties[key] = property
      } else {
         log.Warnf("mergeProperty: merge collision: %s", destination.Name)
      }
   }
}

// SplunkTypeDefinitions Recursively travel down the astType
func (d *Document) SplunkTypeDefinitions(astType *ast.Type, gqlSchema *ast.SchemaDocument) {
   // log.Printf("SplunkDefinitions: astType.NamedType: %v", astType.NamedType)
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

   // take care of types included in GraphQL Union type
   if def.Kind == ast.Union {
      d.SplunkUnionTypeDefinitions(astType, gqlSchema)
   }

   // for all fields for that type
   for _, field := range def.Fields {
      fieldTypeName := nameFromType(field.Type)

      fieldTypeDef := gqlSchema.Definitions.ForName(fieldTypeName)

      if fieldTypeDef == nil { // if it's not a created type then it's not defined in the schema document, it's a basic type
         fieldTypeDef = NewBasicTypeDefinition(field.Name, field.Directives)
      }

      fieldType := field.Type

      fieldTypeProperty, err := NewProperty(fieldTypeDef, fieldType)
      if err != nil { // if not any type, past basic type, then will give error and continue to next field for that type
         log.Warnf("error getting field type as a property: %v", err)
         continue
      }
      if fieldTypeProperty == nil {
         continue
      }

      property.Properties[uppercaseTypeName(field.Name)] = fieldTypeProperty // add property types to this larger property

      if field.Type.NonNull { // if it's required, add name to Required for this property
         property.Required = append(property.Required, uppercaseTypeName(field.Name))
      }

      // check if it's a known type to avoid duplicates/infinite recursion
      // log.Printf("addType: %v", field.Type.Name())
      if err = d.addType(field.Type.Name(), true); err != nil {
         continue
      }

      d.SplunkTypeDefinitions(field.Type, gqlSchema) // call recursively again for the types under the current type

   }
   err = d.AddDefinition(astType, property) // add to the definitions on the doc, that type and its properties
   if err != nil {                          // check duplicates, if so then wasn't added again because already defined
      // log.Warnf("error adding definition to document: %v", err)
   }

}

// SplunkUnionTypeDefinitions
// Run through types under union
func (d *Document) SplunkUnionTypeDefinitions(astType *ast.Type, gqlSchema *ast.SchemaDocument) {
   var def *ast.Definition
   def = gqlSchema.Definitions.ForName(astType.Name())
   if def == nil {
      return
   }

   // log.Printf("%v", def)

   // for each type in GraphQL union type
   for _, unionType := range def.Types {
      fieldTypeDef := gqlSchema.Definitions.ForName(unionType)
      log.Printf("%v", fieldTypeDef.Interfaces)
      // create new ast.Type with type definition
      astFieldType := ast.Type{
         NamedType: fieldTypeDef.Name,
         Elem:      nil,
         NonNull:   false,
         Position:  fieldTypeDef.Position,
      }
      d.SplunkTypeDefinitions(&astFieldType, gqlSchema) // recurse on this type
   }

}

// map of known types to check duplicates

func (d *Document) addType(name string, add bool) (err error) {
   // Remove not null indicator
   name = strings.TrimSuffix(name, "!")

   if _, found := d.knownTypes[name]; found {
      return fmt.Errorf("duplicate type: %v", name)
   }
   // if add is true, modify the map, otherwise only checking the map
   if add {
      d.knownTypes[name] = nil
   }
   return
}
