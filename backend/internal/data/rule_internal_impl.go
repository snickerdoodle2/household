package data

import (
	"encoding/json"
	"fmt"
	"inzynierka/internal/data/validator"
	"slices"
	"time"

	"github.com/google/uuid"
)

type RuleAnd struct {
	Children []RuleInternal `json:"children"`
}

type FakeAnd RuleAnd

func (r RuleAnd) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FakeAnd
		Type string `json:"type"`
	}{
		FakeAnd: FakeAnd(r),
		Type:    "and",
	})
}

func (r *RuleAnd) Process(data RuleData, m *SensorMeasurementModel) (bool, error) {
	for _, child := range r.Children {
		ret, err := child.Process(data, m)
		if err != nil {
			return false, err
		}

		if !ret {
			return false, nil
		}
	}

	return true, nil
}

func (r *RuleAnd) Dependencies() []uuid.UUID {
	res := make([]uuid.UUID, 0)

	for _, child := range r.Children {
		for _, dep := range child.Dependencies() {
			if !slices.Contains(res, dep) {
				res = append(res, dep)
			}
		}
	}

	return res
}

func (r *RuleAnd) Validate(v *validator.Validator) {
	v.Check(len(r.Children) < 1, "and", "must have at least one child")
	for _, child := range r.Children {
		child.Validate(v)
	}
}

type RuleGT struct {
	SensorID uuid.UUID `json:"sensor_id"`
	Value    float64   `json:"value"`
}

type FakeGT RuleGT

func (r RuleGT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FakeGT
		Type string `json:"type"`
	}{
		FakeGT: FakeGT(r),
		Type:   "gt",
	})
}

func (r *RuleGT) Process(data RuleData, _ *SensorMeasurementModel) (bool, error) {
	val, ok := data[r.SensorID]

	if !ok {
		return false, ErrMissingVal
	}

	return val > r.Value, nil
}

func (r *RuleGT) Dependencies() []uuid.UUID {
	return []uuid.UUID{r.SensorID}
}

func (r *RuleGT) Validate(v *validator.Validator) {
}

type RuleLT struct {
	SensorID uuid.UUID `json:"sensor_id"`
	Value    float64   `json:"value"`
}

type FakeLT RuleLT

func (r RuleLT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FakeLT
		Type string `json:"type"`
	}{
		FakeLT: FakeLT(r),
		Type:   "lt",
	})
}

func (r *RuleLT) Process(data RuleData, _ *SensorMeasurementModel) (bool, error) {
	val, ok := data[r.SensorID]

	if !ok {
		return false, ErrMissingVal
	}

	return val < r.Value, nil
}

func (r *RuleLT) Dependencies() []uuid.UUID {
	return []uuid.UUID{r.SensorID}
}

func (r *RuleLT) Validate(v *validator.Validator) {
}

type RuleNot struct {
	Wrapped RuleInternal `json:"wrapped"`
}

type FakeNot RuleNot

func (r RuleNot) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FakeNot
		Type string `json:"type"`
	}{
		FakeNot: FakeNot(r),
		Type:    "not",
	})
}

func (r *RuleNot) Process(data RuleData, m *SensorMeasurementModel) (bool, error) {
	val, err := r.Wrapped.Process(data, m)
	return !val, err
}

func (r *RuleNot) Dependencies() []uuid.UUID {
	return r.Wrapped.Dependencies()
}

func (r *RuleNot) Validate(v *validator.Validator) {
}

type RuleOr struct {
	Children []RuleInternal `json:"children"`
}

type FakeOr RuleOr

func (r RuleOr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FakeOr
		Type string `json:"type"`
	}{
		FakeOr: FakeOr(r),
		Type:   "or",
	})
}

func (r *RuleOr) Process(data RuleData, m *SensorMeasurementModel) (bool, error) {
	for _, child := range r.Children {
		ret, err := child.Process(data, m)
		if err != nil {
			return false, err
		}

		if ret {
			return true, nil
		}
	}

	return false, nil
}

func (r *RuleOr) Dependencies() []uuid.UUID {
	res := make([]uuid.UUID, 0)

	for _, child := range r.Children {
		for _, dep := range child.Dependencies() {
			if !slices.Contains(res, dep) {
				res = append(res, dep)
			}
		}
	}

	return res
}

func (r *RuleOr) Validate(v *validator.Validator) {
	v.Check(len(r.Children) < 1, "or", "must have at least one child")
	for _, child := range r.Children {
		child.Validate(v)
	}
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Duration(d).String())), nil
}

type RulePerc struct {
	SensorID   uuid.UUID `json:"sensor_id"`
	Percentile int       `json:"perc"`
	Delta      Duration  `json:"duration"`
}

type FakePerc RulePerc

func (r RulePerc) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		FakePerc
		Type string `json:"type"`
	}{
		FakePerc: FakePerc(r),
		Type:     "perc",
	})
}

func (r *RulePerc) Process(data RuleData, m *SensorMeasurementModel) (bool, error) {
	val, ok := data[r.SensorID]

	if !ok {
		return false, ErrMissingVal
	}

	perc, err := m.GetPercentile(r.SensorID, time.Duration(r.Delta), r.Percentile)
	if err != nil {
		return false, err
	}
	return val >= perc, nil
}

func (r *RulePerc) Dependencies() []uuid.UUID {
	return []uuid.UUID{r.SensorID}
}

func (r *RulePerc) Validate(v *validator.Validator) {
	v.Check(r.Percentile > 0, "rulePerc", "Percentile should be larger than 0")
	v.Check(r.Percentile <= 0, "rulePerc", "Percentile smaller or equal 100")
	v.Check(r.Delta > 0, "rulePerc", "Duration should be larger than 0")
}
