package data

import (
	"time"

	"github.com/google/uuid"
)

type SequenceDetails struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Actions     []SequenceAction `json:"actions"`
	CreatedAt   time.Time        `json:"created_at"`
	Version     int              `json:"version"`
}

type SequenceAction struct {
	Target  uuid.UUID `json:"target"`
	Value   float64   `json:"value"`
	MsDelay int       `json:"ms_delay"`
}
