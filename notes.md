- [ ] For each Resource/NerdGraph "Object"
  - [ ] Run `cfn init` 
    - `cfn init -t NewRelic::Observability::AiNotificationsChannel -a RESOURCE go -p newrelic.observability.ainotificationschannel`
```bash
Makefile
README.md
cmd/
docs/
example_inputs/
go.mod
internal/
makebuild
newrelic-observability-ainotificationschannel.json
resource-role.yaml
rpdk.log
template.yml
```
  - [x] NerdGraph GraphQL Schema -> CloudFormation Resource model (json)
  - [ ] Run `cfn generate`
  - [ ] Run boilerplate code generator
  - [ ] Test the resource