package rule

import (
	"encoding/json"
	"inzynierka/internal/data/validator"

	"github.com/google/uuid"
)

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

func (r *RuleGT) Validate(v *validator.Validator) {
}
