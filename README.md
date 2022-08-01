Translate a GraphQL Schema document into the corresponding CloudFormation JSON Schema definition
- Read the GraphQL Schema document into an AST
- Recursively traverse the AST building up a Go object representation of the AWS JSON Schema definitions
- Dump the Go representation tree to JSON

Helpful links
- [CloudFormation JSON Schema definitions](https://github.com/aws-cloudformation/cloudformation-cli/tree/master/src/rpdk/core/data/schema)
- [NerdGraph GraphQL Schema definition](schema.graphql)

Build & run
- `go build cmd/gqlparser/main.go  ; ./main > main.json`