package rule

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// TODO: Add stack trace / sensor id to this erorr value
	ErrMissingVal = errors.New("Missing value in data")
)

type RuleData map[uuid.UUID]float64

type RuleInternal interface {
	Process(data RuleData) (bool, error)
	Dependencies() []uuid.UUID
}

type RuleGT struct {
	SensorID uuid.UUID
	Value    float64
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
