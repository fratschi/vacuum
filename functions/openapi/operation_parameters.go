// Copyright 2020-2021 Dave Shanley / Quobix
// SPDX-License-Identifier: MIT

package openapi

import (
	"fmt"
	"github.com/daveshanley/vacuum/model"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
	"strings"
)

// OperationParameters is a rule that checks for valid parameters and parameters combinations
type OperationParameters struct {
}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the OperationParameters rule.
func (op OperationParameters) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{
		Name: "operation_parameters",
	}
}

// RunRule will execute the OperationParameters rule, based on supplied context and a supplied []*yaml.Node slice.
func (op OperationParameters) RunRule(nodes []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	if len(nodes) <= 0 {
		return nil
	}

	var results []model.RuleFunctionResult

	// add any param indexing errors already found.
	errs := context.Index.GetOperationParametersIndexErrors()
	for n := range errs {
		er := errs[n].(*index.IndexingError)
		results = append(results, model.RuleFunctionResult{
			Message:   er.Error(),
			StartNode: er.Node,
			EndNode:   er.Node,
			Path:      er.Path,
			Rule:      context.Rule,
		})
	}

	// look in the index for all operations params.
	for path, methods := range context.Index.GetOperationParameterReferences() {
		for method, methodNode := range methods {

			seenParamInLocations := make(map[string]bool)

			currentVerb := method
			currentPath := path

			resultPath := fmt.Sprintf("$.paths.%s.%s.parameters", path, currentVerb)

			for key, params := range methodNode {

				// TODO: come back and re-visit this code
				if strings.Contains(key, "~1") {
					results = append(results, model.RuleFunctionResult{
						Message: fmt.Sprintf("There is a `~1` character in this `%s` operation at '%s",
							currentVerb, currentPath),
						StartNode: nil,
						EndNode:   nil,
						Path:      resultPath,
						Rule:      context.Rule,
					})
					continue
				}

				for _, param := range params {
					_, paramInNode := utils.FindKeyNode("in", param.Node.Content)
					startNode := param.Node
					endNode := utils.FindLastChildNodeWithLevel(startNode, 0)

					if paramInNode != nil {
						if seenParamInLocations[paramInNode.Value] {
							if paramInNode.Value == "body" {
								results = append(results, model.RuleFunctionResult{
									Message: fmt.Sprintf("the `%s` operation at path `%s` contains a "+
										"duplicate param in:body definition", currentVerb, currentPath),
									StartNode: startNode,
									EndNode:   endNode,
									Path:      resultPath,
									Rule:      context.Rule,
								})
							}
						} else {
							if paramInNode.Value == "body" || paramInNode.Value == "formData" {
								if seenParamInLocations["formData"] || seenParamInLocations["body"] {
									results = append(results, model.RuleFunctionResult{
										Message: fmt.Sprintf("the `%s` operation at path `%s` "+
											"contains parameters using both in:body and in:formData",
											currentVerb, currentPath),
										StartNode: startNode,
										EndNode:   endNode,
										Path:      resultPath,
										Rule:      context.Rule,
									})
								}
							}
							seenParamInLocations[paramInNode.Value] = true
						}
					} else {
						rfr := model.RuleFunctionResult{
							Message: fmt.Sprintf("the `%s` operation at path `%s` contains a "+
								"parameter with no `in` value", currentVerb, currentPath),
							StartNode: startNode,
							EndNode:   endNode,
							Path:      resultPath,
							Rule:      context.Rule,
						}
						results = append(results, rfr)

					}
				}
			}
		}
	}

	return results
}
