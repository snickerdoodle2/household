package rule_test

import (
	"encoding/json"
	"inzynierka/internal/data/rule"
	"slices"
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

func TestRuleUnmarshalling(t *testing.T) {
	data := `{
    "id": "78c6e961-b410-4130-8254-50257b3d88f5",
    "description": "Nowa reguła",
    "internal": {
        "type": "or",
        "children": [
            {
                "type": "and",
                "children": [
                    {
                        "type": "lt",
                        "sensor_id": "7b55654c-fbd1-4054-9b93-228e8e7e8544",
                        "value": 5
                    },
                    {
                        "type": "gt",
                        "sensor_id": "7b55654c-fbd1-4054-9b93-228e8e7e8544",
                        "value": 4
                    }
                ]
            },
            {
                "wrapped": {
                    "sensor_id": "cfe7987c-5ca8-4ad1-8c1e-507ea937d71e",
                    "value": 8,
                    "type": "gt"
                },
                "type": "not"
            }

        ]
    },
    "on_valid": {
        "to": "3a415307-7845-4f05-a790-4e8e203a49c3",
        "payload": {
            "data": "loool"
        }
    }
}`

	a := rule.Rule{}
	err := json.Unmarshal([]byte(data), &a)
	if err != nil {
		t.Fatalf("Expected success, found %v", err)
	}

	validDeps := []string{"7b55654c-fbd1-4054-9b93-228e8e7e8544", "cfe7987c-5ca8-4ad1-8c1e-507ea937d71e"}

	deps := a.Internal.Dependencies()

	for _, valid := range validDeps {
		uuid := uuid.MustParse(valid)
		if !slices.Contains(deps, uuid) {
			t.Errorf("Slice does not contain %v", uuid)
		}
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	tmp := uuid.New()
	internal := rule.RuleOr{
		Children: []rule.RuleInternal{
			&rule.RuleAnd{
				Children: []rule.RuleInternal{
					&rule.RuleLT{
						SensorID: tmp,
						Value:    5,
					},
					&rule.RuleGT{
						SensorID: tmp,
						Value:    8,
					},
				},
			},
			&rule.RuleNot{
				Wrapped: &rule.RuleGT{
					SensorID: uuid.New(),
					Value:    8,
				},
			},
		},
	}
	iRule := rule.Rule{
		ID:          uuid.New(),
		Description: "Nowa reguła",
		Internal:    &internal,
		OnValid: rule.ValidRuleAction{
			To:      uuid.New(),
			Payload: map[string]string{"data": "loool"},
		},
	}

	marshalled, err := json.Marshal(iRule)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	uRule := rule.Rule{}
	err = json.Unmarshal(marshalled, &uRule)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if uRule.ID != iRule.ID {
		t.Errorf("Expected: %v; Got: %v", iRule.ID, uRule.ID)
	}

	if uRule.Description != iRule.Description {
		t.Errorf("Expected: %v; Got: %v", iRule.Description, uRule.Description)
	}

	iDeps := iRule.Internal.Dependencies()
	uDeps := uRule.Internal.Dependencies()

	for _, child := range iDeps {
		if !slices.Contains(uDeps, child) {
			t.Errorf("Missing %v deps", child)
		}
	}
}
