package rule

import (
	"encoding/json"
	"inzynierka/internal/data/validator"

	"github.com/google/uuid"
)

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

func (r *RuleLT) Validate(v *validator.Validator) {
}
