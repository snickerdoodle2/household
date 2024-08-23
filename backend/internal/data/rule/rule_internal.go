package rule

import (
	"encoding/json"
	"errors"
	"slices"

	"github.com/google/uuid"
)

var (
	// TODO: Add stack trace / sensor id to this erorr value
	ErrMissingVal = errors.New("Missing value in data")
)

type RuleData map[uuid.UUID]float64

type RuleInternal interface {
	Process(data RuleData) (bool, error)
	// NOTE: map[uuid.UUID]struct{} (hashset) -> better perf
	// should be called only once per rule lifetime
	Dependencies() []uuid.UUID
}

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

func (r *RuleAnd) Process(data RuleData) (bool, error) {
	for _, child := range r.Children {
		ret, err := child.Process(data)
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
	res := make([]uuid.UUID, len(r.Children))

	for _, child := range r.Children {
		for _, dep := range child.Dependencies() {
			if !slices.Contains(res, dep) {
				res = append(res, dep)
			}
		}
	}

	return res
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

func (r *RuleOr) Process(data RuleData) (bool, error) {
	for _, child := range r.Children {
		ret, err := child.Process(data)
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
	res := make([]uuid.UUID, len(r.Children))

	for _, child := range r.Children {
		for _, dep := range child.Dependencies() {
			if !slices.Contains(res, dep) {
				res = append(res, dep)
			}
		}
	}

	return res
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

func (r *RuleGT) Process(data RuleData) (bool, error) {
	val, ok := data[r.SensorID]

	if !ok {
		return false, ErrMissingVal
	}

	return val > r.Value, nil
}

func (r *RuleGT) Dependencies() []uuid.UUID {
	return []uuid.UUID{r.SensorID}
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

func (r *RuleLT) Process(data RuleData) (bool, error) {
	val, ok := data[r.SensorID]

	if !ok {
		return false, ErrMissingVal
	}

	return val < r.Value, nil
}

func (r *RuleLT) Dependencies() []uuid.UUID {
	return []uuid.UUID{r.SensorID}
}
