package main

import (
   "GraphQLSchema-to-CloudFormationSchema/pkg/aws/cloudformation/model"
   "encoding/json"
   "fmt"
   log "github.com/sirupsen/logrus"
   "github.com/vektah/gqlparser/v2/ast"
   "github.com/vektah/gqlparser/v2/parser"
   "io/ioutil"
)

func main() {
   // Load the GraphQL Schema into an AST
   source, err := ioutil.ReadFile("schema.graphql")
   if err != nil {
      log.Fatalf("error reading schema.json file: %v", err)
   }

   src := ast.Source{
      Name:    "",
      Input:   string(source),
      BuiltIn: true,
   }

   document, gqlerr := parser.ParseSchema(&src)
   if gqlerr != nil {
      log.Fatalf("error parsing schema.json: %v", gqlerr)
   }

   jsonDocument := model.NewDocument()

   for _, mutation := range getMutationTypes(document) {
      mutationDefinition := document.Definitions.ForName(mutation)
      // TODO make the mutation name(s) a command line parameter
      targetMutation := mutationDefinition.Fields.ForName("dashboardCreate")
      log.Printf("main: targetMutation: Name: %+v Type: %v Arguments: %+v", targetMutation.Name, targetMutation.Type, targetMutation.Arguments)
      // NOTE: args go in properties!
      for _, argDef := range targetMutation.Arguments {
         argDef.Description = ""
         log.Printf("main: argDef: %+v", argDef)
         def := document.Definitions.ForName(argDef.Type.NamedType)
         if def == nil {
            // Then it's a schema built-in type like Int, covert the argDef to a vanilla Definition
            def = model.NewBasicTypeDefinition(argDef.Name, argDef.Directives)
         }
         property, err := model.NewProperty(def, argDef.Type)
         if err != nil {
            log.Errorf("error processing Defintion: %v", err)
            continue
         }
         jsonDocument.AddProperty(argDef, property.AsSchemaProperty())
         jsonDocument.SplunkTypeDefinitions(argDef.Type, document)
      }
   }
   log.Printf("main: complete")
   b, err := json.MarshalIndent(jsonDocument, "", "   ")
   if err != nil {
      log.Fatal(err)
   }

   fmt.Println(string(b))

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
