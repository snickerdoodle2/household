package rule

import (
	"encoding/json"
	"inzynierka/internal/data/validator"

	"github.com/google/uuid"
)

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

func (r *RuleNot) Process(data RuleData) (bool, error) {
	val, err := r.Wrapped.Process(data)
	return !val, err
}

func (r *RuleNot) Dependencies() []uuid.UUID {
	return r.Wrapped.Dependencies()
}

func (r *RuleNot) Validate(v *validator.Validator) {
}
