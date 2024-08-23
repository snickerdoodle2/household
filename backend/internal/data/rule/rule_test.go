package rule_test

import (
	"encoding/json"
	"inzynierka/internal/data/rule"
	"testing"

	"github.com/google/uuid"
)

func TestRuleMarshallingNoError(t *testing.T) {
	internalOR := rule.RuleOr{
		Children: []rule.RuleInternal{
			&rule.RuleLT{
				SensorID: uuid.New(),
				Value:    5,
			},
			&rule.RuleGT{
				SensorID: uuid.New(),
				Value:    8,
			},
		},
	}
	rule := rule.Rule{
		ID:          uuid.New(),
		Description: "Nowa reguła",
		Internal:    &internalOR,
		OnValid: rule.ValidRuleAction{
			To:      uuid.New(),
			Payload: map[string]string{"data": "loool"},
		},
	}

	data, err := json.MarshalIndent(rule, "", "    ")

	if err != nil {
		t.Error("Expected success...")
	}

	t.Logf("%s", string(data))
}
