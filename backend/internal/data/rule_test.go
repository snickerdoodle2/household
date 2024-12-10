package data_test

import (
	"encoding/json"
	"inzynierka/internal/data"
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestRuleMarshallingNoError(t *testing.T) {
	internalOR := data.RuleOr{
		Children: []data.RuleInternal{
			&data.RuleLT{
				SensorID: uuid.New(),
				Value:    5,
			},
			&data.RuleGT{
				SensorID: uuid.New(),
				Value:    8,
			},
		},
	}
	rule := data.Rule{
		ID:          uuid.New(),
		Description: "Nowa reguła",
		Internal:    &internalOR,
		OnValid: data.ValidRuleAction{
			TargetType: data.SensorTarget,
			TargetId:   uuid.New(),
			Payload:    map[string]interface{}{"data": "loool"},
		},
	}

	data, err := json.MarshalIndent(rule, "", "    ")

	if err != nil {
		t.Error("Expected success...")
	}

	t.Logf("%s", string(data))
}

func TestRuleUnmarshalling(t *testing.T) {
	jsonData := `{
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
		"target-type": "sensor",
        "target-id": "3a415307-7845-4f05-a790-4e8e203a49c3",
        "payload": {
            "data": "loool"
        }
    }
}`

	a := data.Rule{}
	err := json.Unmarshal([]byte(jsonData), &a)
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
	internal := data.RuleOr{
		Children: []data.RuleInternal{
			&data.RuleAnd{
				Children: []data.RuleInternal{
					&data.RuleLT{
						SensorID: tmp,
						Value:    5,
					},
					&data.RuleGT{
						SensorID: tmp,
						Value:    8,
					},
				},
			},
			&data.RuleNot{
				Wrapped: &data.RuleGT{
					SensorID: uuid.New(),
					Value:    8,
				},
			},
		},
	}
	iRule := data.Rule{
		ID:          uuid.New(),
		Name:        "Nowa reguła",
		Description: "Przykładowy opis nowej reguły",
		Internal:    &internal,
		OnValid: data.ValidRuleAction{
			TargetType: data.SensorTarget,
			TargetId:   uuid.New(),
			Payload:    map[string]interface{}{"data": "loool"},
		},
	}

	marshalled, err := json.Marshal(iRule)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	uRule := data.Rule{}
	err = json.Unmarshal(marshalled, &uRule)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if uRule.ID != iRule.ID {
		t.Errorf("Expected: %v; Got: %v", iRule.ID, uRule.ID)
	}

	if uRule.Name != iRule.Name {
		t.Errorf("Expected: %v; Got: %v", iRule.Name, uRule.Name)
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
