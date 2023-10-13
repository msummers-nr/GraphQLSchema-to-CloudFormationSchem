package nerdgraph

import (
   "GraphQLSchema-to-CloudFormationSchema/pkg/aws/cloudformation/model"
   "encoding/json"
   log "github.com/sirupsen/logrus"
   "github.com/vektah/gqlparser/v2/ast"
   "os"
   "strings"
)

type Service struct {
   serviceName      string
   createDefinition *ast.FieldDefinition
   updateDefinition *ast.FieldDefinition
   deleteDefinition *ast.FieldDefinition
   schemaDocument   *ast.SchemaDocument
}

var services = make(map[string]*Service)

func NewService(definition *ast.FieldDefinition, document *ast.SchemaDocument) *Service {
   serviceName := ParseServiceName(definition.Name)
   service := services[serviceName]
   if service == nil {
      service = &Service{serviceName: serviceName}
      services[serviceName] = service
   }

   // Insert the Tags property
   entity := document.Definitions.ForName("Entity")
   if entity != nil {
      tags := entity.Fields.ForName("tags")
      if tags != nil {
         tagArg := ast.ArgumentDefinition{
            Description:              "",
            Name:                     "tags",
            DefaultValue:             nil,
            Type:                     tags.Type,
            Directives:               nil,
            Position:                 nil,
            BeforeDescriptionComment: nil,
            AfterDescriptionComment:  nil,
         }
         definition.Arguments = append(definition.Arguments, &tagArg)
      }
   }

   if strings.Contains(definition.Name, "Create") {
      service.createDefinition = definition
   } else if strings.Contains(definition.Name, "Update") {
      service.updateDefinition = definition
   } else if strings.Contains(definition.Name, "Delete") {
      service.deleteDefinition = definition
   } else {
      log.Warnf("NewService: ignoring unknown mutation type: %s", definition.Name)
      return nil
   }
   service.schemaDocument = document
   return service
}

func ParseServiceName(name string) string {
   name = strings.ReplaceAll(name, "Create", "")
   name = strings.ReplaceAll(name, "Update", "")
   name = strings.ReplaceAll(name, "Delete", "")
   return name
}

func (s *Service) Emit() {
   doc := model.NewDocument()
   doc.TypeName = "NewRelic::Observability::" + s.serviceName

   // Create GOES First!
   if s.createDefinition != nil {
      s.parse(s.createDefinition, doc)
   }
   if s.updateDefinition != nil {
      s.parse(s.updateDefinition, doc)
   }
   if s.deleteDefinition != nil {
      s.parse(s.deleteDefinition, doc)
   }

   // required must contain unique values
   m := make(map[string]string)
   for _, r := range doc.Required {
      m[r] = ""
   }
   doc.Required = make([]string, 0, len(m))
   for k, _ := range m {
      doc.Required = append(doc.Required, k)
   }

   s.toFile(doc)
}

func (s *Service) GetName() string {
   return s.serviceName
}

func (s *Service) toFile(doc *model.Document) {
   b, err := json.MarshalIndent(doc, "", "   ")
   if err != nil {
      log.Errorf("toFile: error: %v", err)
   }

   fileName := strings.ToLower(doc.TypeName)
   fileName = strings.ReplaceAll(fileName, "::", "-")
   f, err := os.Create(fileName + ".json")
   if err != nil {
      log.Errorf("toFile: error: %v", err)
   }
   defer f.Close()

   _, err = f.WriteString(string(b))
   if err != nil {
      log.Errorf("toFile: error: %v", err)
   }
}

func (s *Service) parse(field *ast.FieldDefinition, doc *model.Document) {
   if field.Arguments == nil {
      recurseFieldTypes(s.schemaDocument, doc, field)
   } else {
      recurseArgTypes(s.schemaDocument, doc, field)
   }
}
