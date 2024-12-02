package data

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SequenceModel struct {
	DB *pgxpool.Pool
}

type Sequence struct {
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
	MsDelay int       `json:"msDelay"`
}

func (m SequenceModel) Insert(sequence *Sequence) error {
	query := `INSERT INTO SEQUENCES (id, name, description, actions)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at, version`

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	sequence.ID = id

	args := []any{sequence.ID, sequence.Name, sequence.Description, sequence.Actions}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = m.DB.QueryRow(ctx, query, args...).Scan(&sequence.CreatedAt, &sequence.Version)

	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint \"uri_unique\""):
			return ErrDuplicateUri
		default:
			return err
		}
	}

	return nil
}
