package main

import (
	"GraphQLSchema-to-CloudFormationSchema/pkg/aws/cloudformation/model"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"io/ioutil"
	"os"
)

func main() {
	// Load the GraphQL Schema into an AST
	var schema string
	fmt.Printf("Enter GraphQL Schema to translate: ")
	fmt.Scan(&schema)
	source, err := ioutil.ReadFile(schema) //ex: schema.graphql
	if err != nil {
		log.Fatalf("error reading file %v: %v", schema, err)
	}

	src := ast.Source{
		Name:    "",
		Input:   string(source),
		BuiltIn: true,
	}

	// Parse the document
	document, gqlerr := parser.ParseSchema(&src)
	if gqlerr != nil {
		log.Fatalf("error parsing document: %v", gqlerr)
	}

	// Create a new model for mutation type output
	jsonDocumentMutation := model.NewDocument()

	// Create a new model for query type output
	jsonDocumentQuery := model.NewDocument()

	//For GraphQL mutation types
	for _, mutation := range getMutationTypes(document) {
		mutationDefinition := document.Definitions.ForName(mutation)

		var targetMutationName string
		fmt.Printf("Enter a target mutation type to translate or press 1 to exit: ")
		fmt.Scan(&targetMutationName) //test - dashboardCreate
		if targetMutationName == "1" {
			break
		}
		targetMutation := mutationDefinition.Fields.ForName(targetMutationName)

		log.Printf("main: targetMutation: Name: %+v Type: %v Arguments: %+v", targetMutation.Name, targetMutation.Type, targetMutation.Arguments)

		//handle recursion through types based on if definition arguments exist
		if targetMutation.Arguments == nil {
			recurseFieldTypes(document, jsonDocumentMutation, targetMutation)
		} else {
			// Uncomment below to also run through targetMutation.Type, not required, warning: increases schema size significantly
			//jsonDocumentMutation.SplunkTypeDefinitions(targetMutation.Type, document)
			recurseArgTypes(document, jsonDocumentMutation, targetMutation)
		}
	}

	//For GraphQL query types
	for _, query := range getQueryTypes(document) {
		queryDefinition := document.Definitions.ForName(query)

		var targetQueryName string
		fmt.Printf("Enter a target query type to translate or press 1 to exit: ")
		fmt.Scan(&targetQueryName)
		if targetQueryName == "1" {
			break
		}
		targetQuery := queryDefinition.Fields.ForName(targetQueryName)

		//handle recursion through types based on if definition arguments exist
		if targetQuery.Arguments == nil {
			recurseFieldTypes(document, jsonDocumentQuery, targetQuery)
		} else {
			recurseArgTypes(document, jsonDocumentQuery, targetQuery)
		}

	}

	log.Printf("main: complete")
	b, err := json.MarshalIndent(jsonDocumentMutation, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	q, err2 := json.MarshalIndent(jsonDocumentQuery, "", "   ")
	if err2 != nil {
		log.Fatal(err2)
	}

	//print string to console
	//fmt.Println(string(b), string(q))

	// print translated schema to its own json file
	writeToFile(string(b), "translated-mutation-schema.json")
	writeToFile(string(q), "translated-query-schema.json")
	fmt.Println("Refer to translated-x-schema.json for translated schema.")
}

func getMutationTypes(document *ast.SchemaDocument) []string {
	types := make([]string, 0)
	for _, schemaDefinition := range document.Schema {
		for _, operationType := range schemaDefinition.OperationTypes {
			if operationType.Operation == "mutation" {
				types = append(types, operationType.Type)
			}
		}
	}
	return types
}

func getQueryTypes(document *ast.SchemaDocument) []string {
	types := make([]string, 0)
	for _, schemaDefinition := range document.Schema {
		for _, operationType := range schemaDefinition.OperationTypes {
			if operationType.Operation == "query" {
				types = append(types, operationType.Type)
			}
		}
	}
	return types
}

func writeToFile(schema, fileName string) {
	f, errFile := os.Create(fileName)
	if errFile != nil {
		log.Fatal(errFile)
	}
	defer f.Close()
	_, errWrite := f.WriteString(schema)
	if errWrite != nil {
		log.Fatal(errWrite)
	}
}
