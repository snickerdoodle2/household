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
	// NOTE: map[uuid.UUID]struct{} (hashset) -> better perf
	// should be called only once per rule lifetime
	Dependencies() []uuid.UUID
}
