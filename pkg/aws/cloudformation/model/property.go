package model

import (
   "fmt"
   log "github.com/sirupsen/logrus"
   "github.com/vektah/gqlparser/v2/ast"
)

/*
Property is common to the Schema's properties and definitions (https://github.com/aws-cloudformation/cloudformation-cli/blob/865c664b420127d4ff83518d35dfdb3e139a2fa2/src/rpdk/core/data/schema/base.definition.schema.v1.json#L244)
*/
type Property struct {
   InsertionOrder bool   `json:"insertionOrder,omitempty"`
   Title          string `json:"title,omitempty"`
   Description    string `json:"description,omitempty"`
   // `json:"examples,omitempty"`
   // `json:"default,omitempty"`
   // `json:"multipleOf,omitempty"`
   // `json:"maximum,omitempty"`
   // `json:"exclusiveMaximum,omitempty"`
   // `json:"minimum,omitempty"`
   // `json:"exclusiveMinimum,omitempty"`
   // `json:"maxLength,omitempty"`
   // `json:"minLength,omitempty"`
   // `json:"pattern,omitempty"`
   // `json:"items,omitempty"`
   // `json:"maxItems,omitempty"`
   // `json:"minItems,omitempty"`
   // `json:"uniqueItems,omitempty"`
   // `json:"contains,omitempty"`
   // `json:"maxProperties,omitempty"`
   // `json:"minProperties,omitempty"`
   Required             []string             `json:"required,omitempty"`
   Properties           map[string]*Property `json:"properties,omitempty"`
   AdditionalProperties bool                 `json:"additionalProperties,omitempty"`
   // `json:"patternProperties,omitempty"`
   // `json:"dependencies,omitempty"`
   // `json:"const,omitempty"`
   Enum  string `json:"enum,omitempty"`
   Type  string `json:"type,omitempty"`
   Ref   string `json:"$ref,omitempty"`
   Items *Item  `json:"items,omitempty"`
   // `json:"format,omitempty"`
   // `json:"allOf,omitempty"`
   // `json:"anyOf,omitempty"`
   // `json:"oneOf,omitempty"`
   // Below here are for housekeeping, not part of the schema.json
   Name               string             `json:"-"`
   Kind               ast.DefinitionKind `json:"-"`
   BuiltIn            bool               `json:"-"`
   IsRequired         bool               `json:"-"`
   IsArray            bool               `json:"="`
   ArrayEntryRequired bool               `json:"="`
}

type Item struct {
   Type string `json:"type,omitempty"`
   Ref  string `json:"$ref,omitempty"`
}

func NewProperty(definition *ast.Definition, typeDef *ast.Type) (property *Property, err error) {
   property = &Property{
      InsertionOrder:       false,
      Title:                "",
      Description:          "",
      Required:             nil,
      Properties:           make(map[string]*Property),
      AdditionalProperties: false,
      Enum:                 "",
      Ref:                  "",
      Name:                 definition.Name,
      Kind:                 definition.Kind,
      BuiltIn:              definition.BuiltIn,
      IsRequired:           typeDef.NonNull,
      Items:                nil,
   }

   property.IsArray = typeDef.Elem != nil
   name := typeDef.NamedType
   // name := nameFromType(typeDef)
   property.Name = name

   // TODO Refactor this to a func/method- it'll get bigger
   if property.IsArray {
      item := Item{}
      property.Type = "array"
      // TODO is this the correct test?
      if property.Kind == ast.InputObject {
         item.Ref = "#/definitions/" + name
      } else {
         item.Type = name
      }
      property.Items = &item
   } else {
      if property.Kind == ast.InputObject {
         property.Type = ""
         property.Ref = "#/definitions/" + name
      }
   }

   // log.Printf("NewProperty: definition: Name: %v Kind: %v Type: %v BuiltIn: %v Fields: %+v", definition.Name, definition.Kind, definition.Types, definition.BuiltIn, definition.Fields)
   // TODO refactor common to here, then pass *property to the New for the particulars

   switch definition.Kind {
   case ast.Enum:
      err = NewEnum(property, definition, typeDef)
   case ast.Scalar:
      err = NewScalar(property, definition, typeDef)
   case ast.InputObject:
      err = NewInputObject(property, definition, typeDef)
   case ast.Interface:
      err = NewInterface(property, definition, typeDef)
   case ast.Object:
      err = NewObject(property, definition, typeDef)
   case ast.Union:
      err = NewUnion(property, definition, typeDef)
   case "":
      err = NewBasicType(property, definition, typeDef)
   default:
      return nil, fmt.Errorf("unknown Definition.Kind: %v", definition.Kind)
   }
   if err != nil {
      return nil, err
   }
   log.Printf("NewProperty: property: Name: %v Kind: %v TypeDef: %v BuiltIn: %v Ref: %+v IsArray: %v IsRequired: %v ArrayEntryReqquired: %v", property.Name, property.Kind, typeDef, property.BuiltIn, property.Ref, property.IsArray, property.IsRequired, property.ArrayEntryRequired)
   return property, err
}
func NewDefinitionProperty(definition *ast.Definition, typeDef *ast.Type) (property *Property, err error) {
   property = &Property{
      InsertionOrder:       false,
      Title:                "",
      Description:          "",
      Required:             nil,
      Properties:           make(map[string]*Property),
      AdditionalProperties: false,
      Enum:                 "",
      Ref:                  "",
      Name:                 definition.Name,
      Kind:                 definition.Kind,
      BuiltIn:              definition.BuiltIn,
      IsRequired:           typeDef.NonNull,
      Items:                nil,
   }

   property.IsArray = typeDef.Elem != nil
   name := nameFromType(typeDef)
   property.Type = name

   // TODO Refactor this to a func/method- it'll get bigger
   // TODO is this the correct test?
   if property.Kind == ast.InputObject {
      property.Type = "object"
   } else {
      property.Type = name
   }

   // log.Printf("NewProperty: definition: Name: %v Kind: %v Type: %v BuiltIn: %v Fields: %+v", definition.Name, definition.Kind, definition.Types, definition.BuiltIn, definition.Fields)
   // TODO refactor common to here, then pass *property to the New for the particulars

   switch definition.Kind {
   case ast.Enum:
      err = NewEnum(property, definition, typeDef)
   case ast.Scalar:
      err = NewScalar(property, definition, typeDef)
   case ast.InputObject:
      err = NewInputObject(property, definition, typeDef)
   case ast.Interface:
      err = NewInterface(property, definition, typeDef)
   case ast.Object:
      err = NewObject(property, definition, typeDef)
   case ast.Union:
      err = NewUnion(property, definition, typeDef)
   default:
      return nil, fmt.Errorf("unknown Definition.Kind: %v", definition.Kind)
   }
   if err != nil {
      return nil, err
   }
   log.Printf("NewProperty: property: Name: %v Kind: %v TypeDef: %v BuiltIn: %v Ref: %+v IsArray: %v IsRequired: %v ArrayEntryReqquired: %v", property.Name, property.Kind, typeDef, property.BuiltIn, property.Ref, property.IsArray, property.IsRequired, property.ArrayEntryRequired)
   return property, err
}

func propertiesFromFields(fields ast.FieldList, property *Property, properties *[]*Property) {
   for _, field := range fields {
      _ = field
   }
}

func (p *Property) decodeName() {
   // TODO set IsArray
   // TODO set IsRequired
   // TODO set Name to cleaned-up Name
}

func (p *Property) AsSchemaProperty() *Property {
   // TODO massage this property as a JSON Schema top-level "property"
   return p
}
