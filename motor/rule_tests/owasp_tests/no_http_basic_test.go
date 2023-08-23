package tests

import (
	"testing"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
	"github.com/stretchr/testify/assert"
)

func TestRuleSet_OWASPNoHttpBasic_Success(t *testing.T) {

	yml := `openapi: "3.1.0"
info:
  version: "1.0"
components:
  securitySchemes:
    "anything-else":
      type: "http"
      scheme: "bearer"`

	t.Run("valid case", func(t *testing.T) {
		rules := make(map[string]*model.Rule)
		rules["owasp-no-http-basic"] = rulesets.GetOWASPNoHttpBasicRule()

		rs := &rulesets.RuleSet{
			Rules: rules,
		}

		rse := &motor.RuleSetExecution{
			RuleSet: rs,
			Spec:    []byte(yml),
		}
		results := motor.ApplyRulesToRuleSet(rse)
		assert.Len(t, results.Results, 0)
	})
}

func TestRuleSet_OWASPNoHttpBasic_Error(t *testing.T) {

	yml := `openapi: "3.1.0"
info:
  version: "1.0"
components:
  securitySchemes:
    "bad negotiate":
      type: "http"
      scheme: "negotiate"
    "please-hack-me":
      type: "http"
      scheme: basic`

	t.Run("invalid case", func(t *testing.T) {
		rules := make(map[string]*model.Rule)
		rules["owasp-no-http-basic"] = rulesets.GetOWASPNoHttpBasicRule()

		rs := &rulesets.RuleSet{
			Rules: rules,
		}

		rse := &motor.RuleSetExecution{
			RuleSet: rs,
			Spec:    []byte(yml),
		}
		results := motor.ApplyRulesToRuleSet(rse)
		assert.Len(t, results.Results, 1)
	})
}
