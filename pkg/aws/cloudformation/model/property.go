package model

import (
	"fmt"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

/*
Property is common to the Schema's properties and definitions (https://github.com/aws-cloudformation/cloudformation-cli/blob/865c664b420127d4ff83518d35dfdb3e139a2fa2/src/rpdk/core/data/schema/base.definition.schema.v1.json#L244)
*/
type Property struct {
	InsertionOrder *bool  `json:"insertionOrder,omitempty"`
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
	AdditionalProperties *bool                `json:"additionalProperties,omitempty"`
	// `json:"patternProperties,omitempty"`
	// `json:"dependencies,omitempty"`
	// `json:"const,omitempty"`
	Enum  []string `json:"enum,omitempty"`
	Type  string   `json:"type,omitempty"`
	Ref   string   `json:"$ref,omitempty"`
	Items *Item    `json:"items,omitempty"`
	// `json:"format,omitempty"`
	// `json:"allOf,omitempty"`
	AnyOf []*Item `json:"anyOf,omitempty"`
	OneOf []*Item `json:"oneOf,omitempty"`
	// Below here are for housekeeping, not part of the schema.json
	Name               string             `json:"-"`
	Kind               ast.DefinitionKind `json:"-"`
	BuiltIn            bool               `json:"-"`
	IsRequired         bool               `json:"-"`
	IsArray            bool               `json:"="`
	ArrayEntryRequired bool               `json:"="`
}

type Item struct {
	Type        string  `json:"type,omitempty"`
	Ref         string  `json:"$ref,omitempty"`
	AnyOf       []*Item `json:"anyOf,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`
}

// NewProperty
/*
NewProperty create new properties for the schema
*/
func NewProperty(definition *ast.Definition, typeDef *ast.Type) (property *Property, err error) {
	property = &Property{
		Title:       "",
		Description: "",
		Required:    nil,
		Properties:  make(map[string]*Property),
		Enum:        make([]string, 0),
		Ref:         "",
		Name:        definition.Name,
		Kind:        definition.Kind,
		BuiltIn:     definition.BuiltIn,
		IsRequired:  typeDef.NonNull,
		Items:       nil,
	}

	//Handle array
	if IsArrayType(typeDef.String()) {
		property.IsArray = true
	}

	name := nameFromType(typeDef)
	property.Name = name

	//If definition already defined, add $ref instead of writing again; and
	//handle arrays
	property.addRefAndArray()

	// Modify each property based on definition kind
	err = property.createNewDefinitionKind(definition, typeDef)

	if err != nil {
		return nil, err
	}
	//log.Printf("NewProperty: property: Name: %v Kind: %v TypeDef: %v BuiltIn: %v Ref: %+v IsArray: %v IsRequired: %v ArrayEntryRequired: %v", property.Name, property.Kind, typeDef, property.BuiltIn, property.Ref, property.IsArray, property.IsRequired, property.ArrayEntryRequired)
	return property, err
}

// NewDefinitionProperty
/*
NewDefinitionProperty creates a new definition property
*/
func NewDefinitionProperty(definition *ast.Definition, typeDef *ast.Type) (property *Property, err error) {
	property = &Property{
		Title:       "",
		Description: "",
		Required:    nil,
		Properties:  make(map[string]*Property),
		Enum:        make([]string, 0),
		Ref:         "",
		Name:        definition.Name,
		Kind:        definition.Kind,
		BuiltIn:     definition.BuiltIn,
		IsRequired:  typeDef.NonNull,
		Items:       nil,
	}

	//Handle array
	if IsArrayType(typeDef.String()) {
		property.IsArray = true
	}

	name := nameFromType(typeDef)
	property.Name = name

	// Modify each property based on definition kind
	err = property.createNewDefinitionKind(definition, typeDef)

	if err != nil {
		return nil, err
	}

	//log.Printf("NewProperty: property: Name: %v Kind: %v TypeDef: %v BuiltIn: %v Ref: %+v IsArray: %v IsRequired: %v ArrayEntryReqquired: %v", property.Name, property.Kind, typeDef, property.BuiltIn, property.Ref, property.IsArray, property.IsRequired, property.ArrayEntryRequired)
	return property, err
}

func (p *Property) addRefAndArray() {
	if p.IsArray {
		item := Item{}
		// if it's not a basic type
		if p.Kind != "" { //add ref because will definition already added
			item.Ref = "#/definitions/" + p.Name
		} else {
			item.Type = strings.ToLower(p.Name)
		}
		p.Type = "array"
		f := new(bool)
		*f = false
		p.InsertionOrder = f
		p.Items = &item //add to Property Item, property $ref
	} else { //if not array
		if p.Kind == ast.InputObject || p.Kind == ast.Enum || p.Kind == ast.Scalar || p.Kind == ast.Union || p.Kind == ast.Object || p.Kind == ast.Interface {
			p.Type = ""
			p.Ref = "#/definitions/" + p.Name
		}
	}
}

// createNewDefinitionKind
// modify definition property based on Definition.Kind for each Type
func (p *Property) createNewDefinitionKind(definition *ast.Definition, typeDef *ast.Type) (err error) {
	switch definition.Kind {
	case ast.Enum:
		err = NewEnum(p, definition, typeDef)
	case ast.Scalar:
		err = NewScalar(p, definition, typeDef)
	case ast.InputObject:
		err = NewInputObject(p, definition, typeDef)
	case ast.Interface:
		err = NewInterface(p, definition, typeDef)
	case ast.Object:
		err = NewObject(p, definition, typeDef)
	case ast.Union:
		err = NewUnion(p, definition, typeDef)
	case "":
		err = NewBasicType(p, definition, typeDef)
	default:
		return fmt.Errorf("unknown Definition.Kind: %v", definition.Kind)
	}
	return err
}

func (p *Property) AsSchemaProperty() *Property {
	// TODO massage this property as a JSON Schema top-level "property"
	return p
}
