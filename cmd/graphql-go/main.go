package main

import (
   "github.com/graphql-go/graphql/language/ast"
   "github.com/graphql-go/graphql/language/parser"
   "github.com/graphql-go/graphql/language/source"
   "github.com/sirupsen/logrus"
   "io/ioutil"
   "log"
   "reflect"
)

func main() {
   // Load the GraphQL Schema into an AST
   ba, err := ioutil.ReadFile("schema.json.graphql")
   if err != nil {
      logrus.Fatalf("error parsing schema.json: %v", err)
   }
   log.Printf("len(source): %v", len(ba))

   opts := parser.ParseOptions{
      NoLocation: false,
      NoSource:   true,
   }

   src := source.Source{
      Body: ba,
      Name: "schema.json.graphql",
   }
   params := parser.ParseParams{
      Source:  src,
      Options: opts,
   }

   var document *ast.Document
   document, err = parser.Parse(params)
   if err != nil {
      logrus.Fatalf("error parsing schema.json: %v", err)
   }

   log.Printf("%v %+v", reflect.TypeOf(document), document)
   for _, node := range document.Definitions {
      log.Printf("%v", node)
   }
}
