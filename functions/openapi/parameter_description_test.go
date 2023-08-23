package openapi

import (
	"github.com/daveshanley/vacuum/model"
	"github.com/pb33f/libopenapi/index"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestOAS2ParameterDescription_GetSchema(t *testing.T) {
	def := ParameterDescription{}
	assert.Equal(t, "oas2_parameter_description", def.GetSchema().Name)
}

func TestOAS2ParameterDescription_RunRule(t *testing.T) {
	def := ParameterDescription{}
	res := def.RunRule(nil, model.RuleFunctionContext{})
	assert.Len(t, res, 0)
}

func TestOAS2ParameterDescription_RunRule_Success(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      parameters:
        - in: header
          name: blue-eyes
          description: beautiful girl
parameters:
  Maddy:
   in: header
   name: little champion
   description: beautiful boy`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	rule := buildOpenApiTestRuleAction(path, "oas2_parameter_description", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := ParameterDescription{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 0)
}

func TestOAS2ParameterDescription_RunRule_Fail(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      parameters:
        - in: header
          name: blue-eyes
parameters:
  Maddy:
   in: header
   name: little champion`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	rule := buildOpenApiTestRuleAction(path, "oas2_parameter_description", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := ParameterDescription{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 2)
}

func TestOAS2ParameterDescription_RunRule_Success_NoIn(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      parameters:
        - name: blue-eyes
parameters:
  Maddy:
   name: little champion`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	rule := buildOpenApiTestRuleAction(path, "oas2_parameter_description", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := ParameterDescription{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 0)
}

func TestOAS2ParameterDescription_RunRule_Fail_EmptyDescription(t *testing.T) {

	yml := `swagger: 2.0
paths:
  /melody:
    post:
      parameters:
        - in: header
          name: blue-eyes
          description:  
parameters:
  Maddy:
   in: header
   name: little champion
   description:`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	rule := buildOpenApiTestRuleAction(path, "oas2_parameter_description", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := ParameterDescription{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 2)
}

func TestParameterDescription_RunRule_Fail_EmptyDescription_OpenAPI3(t *testing.T) {

	yml := `openapi: 3.0
paths:
  /melody:
    post:
      parameters:
        - in: header
          name: blue-eyes
          description:  
components:
  parameters:
    Maddy:
      in: header
      name: little champion
      description:`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	rule := buildOpenApiTestRuleAction(path, "oas3_parameter_description", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := ParameterDescription{}
	res := def.RunRule(rootNode.Content, ctx)

	assert.Len(t, res, 2)
}
