package rule

import (
	"encoding/json"
	"slices"

	"github.com/google/uuid"
)

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
