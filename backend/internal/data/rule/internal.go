package rule

import "github.com/google/uuid"

type RuleData map[uuid.UUID]float64

type RuleInternal interface {
	Process(data RuleData) (bool, error)
	Dependencies() []uuid.UUID
}
