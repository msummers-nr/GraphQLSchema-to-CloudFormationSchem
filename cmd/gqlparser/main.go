package main

import (
   "GraphQLSchema-to-CloudFormationSchema/pkg/nerdgraph"
   "flag"
   "fmt"
   log "github.com/sirupsen/logrus"
   "github.com/vektah/gqlparser/v2/ast"
   "github.com/vektah/gqlparser/v2/parser"
   "os"
   "strings"
)

func main() {

   // Command line params
   schema := flag.String("schema", "schema.graphql", "File containing the GraphQL Schema to parse")
   list := flag.Bool("list", false, "Set to true to list available mutations and queries")
   mutations := flag.String("mutations", "", "Comma separated list of mutation prefixes to process. Empty == all")
   queries := flag.String("queries", "", "Comma separated list of queries to process. Empty == all")
   logLevel := flag.String("logLevel", "info", "logrus logging level panic | fatal | error | warn | info | debug | trace")
   flag.Parse()

   level, err := log.ParseLevel(*logLevel)
   if err != nil {
      log.Warnf("main: invalid logLevel: %s err: %v", *logLevel, err)
      log.SetLevel(log.InfoLevel)
   } else {
      log.SetLevel(level)
   }
   log.Infof("main: logLevel: %v", log.GetLevel())

   mutationList := strings.Split(*mutations, ",")
   queryList := strings.Split(*queries, ",")
   allMutations := false
   if len(mutationList) == 0 {
      allMutations = true
   }
   allQueries := false
   if len(queryList) == 0 {
      allQueries = true
   }
   _ = allQueries

   // Read the schema file
   source, err := os.ReadFile(*schema) // ex: schema.graphql
   log.Debugf("Reading schema: %s", *schema)
   if err != nil {
      log.Fatalf("error reading file %v: %v", *schema, err)
   }

   // Load the GraphQL Schema into an AST
   src := ast.Source{
      Name:    "",
      Input:   string(source),
      BuiltIn: true,
   }

   // Parse the schemaDocument
   schemaDocument, err := parser.ParseSchema(&src)
   if err != nil {
      log.Fatalf("error parsing schemaDocument: %+v", err)
   }

   services := make(map[string]*nerdgraph.Service)
   for _, mutationDefinition := range getMutationDefinitions(schemaDocument) {
      // At this point we have a HIGH LEVEL (e.g. RootMutationType), the mutations we're interested in are buried in that object's Fields
      for _, fieldDefinition := range mutationDefinition.Fields {
         if *list {
            fmt.Printf("mutation: %s\n", fieldDefinition.Type)
         }

         if allMutations || process(nerdgraph.ParseServiceName(fieldDefinition.Name), mutationList) {
            var service *nerdgraph.Service
            service = nerdgraph.NewService(fieldDefinition, schemaDocument)
            if service != nil {
               services[service.GetName()] = service
            }
         }
      }
   }

   // We've loaded and grouped the services, tell them to parse and marshal
   for _, service := range services {
      service.Emit()
   }
}

func process(mutation string, list []string) bool {
   log.Debugf("process: mutation: %s list: %v", mutation, list)
   for _, prefix := range list {
      if strings.HasPrefix(mutation, prefix) {
         return true
      }
   }
   return false
}

// Get the mutation Definitions associated at the Schema TOP LEVEL.  Most likely this is one- RootMutationType
func getMutationDefinitions(doc *ast.SchemaDocument) []*ast.Definition {
   defs := make([]*ast.Definition, 0, len(doc.Schema))
   for _, sd := range doc.Schema {
      for _, op := range sd.OperationTypes {
         if op.Operation == ast.Mutation {
            defs = append(defs, doc.Definitions.ForName(op.Type))
         }
      }
   }
   return defs
}
