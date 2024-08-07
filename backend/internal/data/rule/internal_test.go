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

func TestRuleAndDependency(t *testing.T) {
	sensorId1 := uuid.New()
	rulegt1 := rule.RuleGT{
		SensorID: sensorId1,
	}

	sensorId2 := uuid.New()
	rulegt2 := rule.RuleGT{
		SensorID: sensorId2,
	}

	ids := []uuid.UUID{sensorId1, sensorId2}

	ruleAnd := rule.RuleAnd{
		Children: []rule.RuleInternal{&rulegt1, &rulegt2},
	}

	deps := ruleAnd.Dependencies()

	for _, id := range ids {
		if !slices.Contains(deps, id) {
			t.Errorf("deps does not contain %q", id)
		}
	}
}

var AndTests = []struct {
	in  []float64
	out bool
}{
	{[]float64{7.5, 10}, true},
	{[]float64{3.5, 10}, false},
	{[]float64{2.5, 8}, false},
}

func TestRuleAndProcess(t *testing.T) {
	sensorId1 := uuid.New()

	rulegt1 := rule.RuleGT{
		SensorID: sensorId1,
		Value:    5,
	}

	sensorId2 := uuid.New()
	rulegt2 := rule.RuleGT{
		SensorID: sensorId2,
		Value:    9,
	}

	ruleAnd := rule.RuleAnd{
		Children: []rule.RuleInternal{&rulegt1, &rulegt2},
	}

	for _, test := range AndTests {
		data := map[uuid.UUID]float64{
			sensorId1: test.in[0],
			sensorId2: test.in[1],
		}

		got, err := ruleAnd.Process(data)
		if err != nil {
			if err != nil {
				t.Errorf("test case %v returned error", test.in)
			}
		}

		if got != test.out {
			t.Errorf("test case %f: wanted %t, got %t", test.in, test.out, got)
		}

	}

}

func TestRuleOrDependency(t *testing.T) {
	sensorId1 := uuid.New()
	rulegt1 := rule.RuleGT{
		SensorID: sensorId1,
	}

	sensorId2 := uuid.New()
	rulegt2 := rule.RuleGT{
		SensorID: sensorId2,
	}

	ids := []uuid.UUID{sensorId1, sensorId2}

	ruleOr := rule.RuleOr{
		Children: []rule.RuleInternal{&rulegt1, &rulegt2},
	}

	deps := ruleOr.Dependencies()

	for _, id := range ids {
		if !slices.Contains(deps, id) {
			t.Errorf("deps does not contain %q", id)
		}
	}
}

var OrTests = []struct {
	in  []float64
	out bool
}{
	{[]float64{7.5, 10}, true},
	{[]float64{3.5, 10}, true},
	{[]float64{2.5, 8}, false},
}

func TestRuleOrProcess(t *testing.T) {
	sensorId1 := uuid.New()

	rulegt1 := rule.RuleGT{
		SensorID: sensorId1,
		Value:    5,
	}

	sensorId2 := uuid.New()
	rulegt2 := rule.RuleGT{
		SensorID: sensorId2,
		Value:    9,
	}

	ruleAnd := rule.RuleOr{
		Children: []rule.RuleInternal{&rulegt1, &rulegt2},
	}

	for _, test := range OrTests {
		data := map[uuid.UUID]float64{
			sensorId1: test.in[0],
			sensorId2: test.in[1],
		}

		got, err := ruleAnd.Process(data)
		if err != nil {
			if err != nil {
				t.Errorf("test case %v returned error", test.in)
			}
		}

		if got != test.out {
			t.Errorf("test case %f: wanted %t, got %t", test.in, test.out, got)
		}

	}

}
