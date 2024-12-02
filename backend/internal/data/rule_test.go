package data_test

import (
	"encoding/json"
	"inzynierka/internal/data"
	"slices"
	"testing"
	"time"

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
			To:      uuid.New(),
			Payload: map[string]interface{}{"data": "loool"},
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
        "to": "3a415307-7845-4f05-a790-4e8e203a49c3",
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
			To:      uuid.New(),
			Payload: map[string]interface{}{"data": "loool"},
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

func TestTimeComp(t *testing.T) {
	cmp := func(hour, minute int, variant data.TimeType, now time.Time) bool {
		switch variant {
		case data.TimeBefore:
			if hour < now.Hour() {
				return true
			}
			if hour == now.Hour() && minute < now.Minute() {
				return true
			}
		case data.TimeAfter:
			if hour > now.Hour() {
				return true
			}
			if hour == now.Hour() && minute > now.Minute() {
				return true
			}
		default:
			panic("Unhandled time rule variant")
		}
		return false
	}

	type Test struct {
		hour     int
		minute   int
		variant  data.TimeType
		now      time.Time
		expected bool
	}

	tests := []Test{
		{
			hour:     23,
			minute:   58,
			variant:  data.TimeBefore,
			now:      time.Date(0, time.January, 0, 23, 59, 0, 0, time.UTC),
			expected: true,
		},
		{
			hour:     23,
			minute:   58,
			variant:  data.TimeBefore,
			now:      time.Date(0, time.January, 0, 23, 58, 0, 0, time.UTC),
			expected: false,
		},
		{
			hour:     21,
			minute:   58,
			variant:  data.TimeBefore,
			now:      time.Date(0, time.January, 0, 22, 58, 0, 0, time.UTC),
			expected: true,
		},
		{
			hour:     21,
			minute:   58,
			variant:  data.TimeAfter,
			now:      time.Date(0, time.January, 0, 22, 58, 0, 0, time.UTC),
			expected: false,
		},
		{
			hour:     21,
			minute:   58,
			variant:  data.TimeAfter,
			now:      time.Date(0, time.January, 0, 20, 58, 0, 0, time.UTC),
			expected: true,
		},
	}

	for i, test := range tests {
		res := cmp(test.hour, test.minute, test.variant, test.now)
		if res != test.expected {
			t.Errorf("Failed test %d: (expected) %t != %t (result)", i, test.expected, res)
		}
	}
}
