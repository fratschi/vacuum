// Copyright 2020-2021 Dave Shanley / Quobix
// SPDX-License-Identifier: MIT

package core

import (
	"fmt"
	"github.com/daveshanley/vacuum/model"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

// Truthy is a rule that will determine if something is seen as 'true' (could be a 1 or "pizza", or actually 'true')
type Truthy struct {
}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the Truthy rule.
func (t Truthy) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{
		Name: "truthy",
	}
}

// RunRule will execute the Truthy rule, based on supplied context and a supplied []*yaml.Node slice.
func (t *Truthy) RunRule(nodes []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	if len(nodes) <= 0 {
		return nil
	}

	var results []model.RuleFunctionResult
	pathValue := "unknown"
	if path, ok := context.Given.(string); ok {
		pathValue = path
	}

	isArray := false
	if len(nodes) == 1 && utils.IsNodeArray(nodes[0]) {
		nodes = nodes[0].Content
		isArray = true
	}

	for x, node := range nodes {

		fieldNode, fieldNodeValue := utils.FindFirstKeyNode(context.RuleAction.Field, node.Content, 0)
		if fieldNode == nil && fieldNodeValue == nil || fieldNodeValue.Value == "false" ||
			fieldNodeValue.Value == "0" || fieldNodeValue.Value == "" {

			if isArray {
				pathValue = fmt.Sprintf("%s[%d]", pathValue, x)
			}

			if !utils.IsNodeMap(fieldNode) && !utils.IsNodeArray(fieldNodeValue) {
				var endNode *yaml.Node
				if len(node.Content) > 0 {
					endNode = node.Content[len(node.Content)-1]
				} else {
					endNode = node
				}
				results = append(results, model.RuleFunctionResult{
					Message:   fmt.Sprintf("%s: `%s` must be set", context.Rule.Description, context.RuleAction.Field),
					StartNode: node,
					EndNode:   endNode,
					Path:      pathValue,
					Rule:      context.Rule,
				})
			}
		}

	}

	return results
}
