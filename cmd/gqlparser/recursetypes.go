package main

import (
	"GraphQLSchema-to-CloudFormationSchema/pkg/aws/cloudformation/model"
	log "github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
)

// recurseFieldTypes
// handle recursion through types if definition only has field type and no arguments
func recurseFieldTypes(document *ast.SchemaDocument, jsonDocument *model.Document, fieldDef *ast.FieldDefinition) {
	def := document.Definitions.ForName(fieldDef.Type.NamedType)
	for _, field := range def.Fields {
		//need for if field.Arguments == nil or != nil, then different (note: if nil, carry on as usual)
		field.Description = ""

		def := handleDefinition(document, field.Type)

		// See if this is a Graphql Schema definition or a builtin
		if def == nil {
			// Then it's a schema built-in type like Int, convert the field to a vanilla Definition
			def = model.NewBasicTypeDefinition(field.Name, field.Directives)
		}

		// Create a new model property
		property, err := model.NewProperty(def, field.Type)
		if err != nil {
			log.Errorf("error processing Definition: %v", err)
			continue
		}
		// Add the property to the output model
		jsonDocument.AddProperty(field.Name, property.AsSchemaProperty())
		// Recursively travel down the field.Type
		jsonDocument.SplunkTypeDefinitions(field.Type, document)
	}
}

// recurseArgTypes
// handle recursion through definition arguments
func recurseArgTypes(document *ast.SchemaDocument, jsonDocument *model.Document, def *ast.FieldDefinition) {
	// NOTE: args go in JSON properties!
	for _, argDef := range def.Arguments { //for each argument under this target mutation
		argDef.Description = ""
		log.Printf("main: argDef: %+v", argDef)

		def := handleDefinition(document, argDef.Type)

		// See if this is a Graphql Schema definition or a builtin
		if def == nil {
			// Then it's a schema built-in type like Int, convert the argDef to a vanilla Definition
			def = model.NewBasicTypeDefinition(argDef.Name, argDef.Directives)
		}

		// Create a new model property
		property, err := model.NewProperty(def, argDef.Type)
		if err != nil {
			log.Errorf("error processing Definition: %v", err)
			continue
		}
		// Add the property to the output model
		jsonDocument.AddProperty(argDef.Name, property.AsSchemaProperty())
		// Recursively travel down the argDef.Type
		jsonDocument.SplunkTypeDefinitions(argDef.Type, document)
	}
}

// handleDefinition
// get definition from GQL schema document, handle array formatting
func handleDefinition(document *ast.SchemaDocument, astType *ast.Type) *ast.Definition {
	var def *ast.Definition
	isArrayProperty := model.IsArrayType(astType.String())
	if isArrayProperty {
		def = document.Definitions.ForName(astType.Name())
	} else {
		def = document.Definitions.ForName(astType.NamedType)
	}
	return def
}
