package rule_test

import (
	"errors"
	"inzynierka/internal/data/rule"
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestRuleGTDependency(t *testing.T) {
	sensorId := uuid.New()

	rule := rule.RuleGT{
		SensorID: sensorId,
		Value:    10,
	}

	deps := rule.Dependencies()

	if len(deps) != 1 {
		t.Errorf("got %d, wanted 1", len(deps))
	}

	if !slices.Contains(deps, sensorId) {
		t.Errorf("returned slice does not contain %q", sensorId)
	}
}

var GTtests = []struct {
	in  float64
	out bool
}{
	{7.5, false},
	{11, true},
	{10, false},
}

func TestRuleGTProcess(t *testing.T) {
	sensorId := uuid.New()

	rule := rule.RuleGT{
		SensorID: sensorId,
		Value:    10,
	}

	for _, test := range GTtests {
		data := map[uuid.UUID]float64{
			sensorId: test.in,
		}

		got, err := rule.Process(data)
		if err != nil {
			t.Errorf("test case %f returned error", test.in)
			continue
		}

		if got != test.out {
			t.Errorf("test case %f: wanted %t, got %t", test.in, test.out, got)
		}
	}

}

func TestRuleGTProcessError(t *testing.T) {
	sensorId := uuid.New()

	rulegt := rule.RuleGT{
		SensorID: sensorId,
		Value:    10,
	}

	data := make(map[uuid.UUID]float64)

	if _, err := rulegt.Process(data); !errors.Is(err, rule.ErrMissingVal) {
		t.Errorf("wanted error, got %s", err.Error())
	}

}
