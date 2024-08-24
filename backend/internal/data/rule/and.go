package rule

import (
	"encoding/json"
	"slices"

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
