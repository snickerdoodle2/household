package data

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	Value   float32   `json:"value"`
	MsDelay int       `json:"msDelay"`
}

type SequenceInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
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

func (m SequenceModel) GetAllInfo() ([]*SequenceInfo, error) {
	query := `SELECT id, name, description
	FROM sequences`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allInfo []*SequenceInfo

	for rows.Next() {
		var info SequenceInfo

		err := rows.Scan(&info.ID, &info.Name, &info.Description)
		if err != nil {
			return nil, err
		}

		allInfo = append(allInfo, &info)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allInfo, nil
}

func (m SequenceModel) GetAll() ([]*Sequence, error) {
	query := `SELECT id, name, description, actions, created_at, version
	FROM sequences`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []*Sequence

	for rows.Next() {
		var sequence Sequence
		var actions []SequenceAction

		err := rows.Scan(
			&sequence.ID,
			&sequence.Name,
			&sequence.Description,
			&actions,
			&sequence.CreatedAt,
			&sequence.Version,
		)
		if err != nil {
			return nil, err
		}

		sequence.Actions = actions
		sequences = append(sequences, &sequence)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sequences, nil
}

func (m SequenceModel) Get(id uuid.UUID) (*Sequence, error) {
	query := `SELECT id, name, description, actions, created_at, version
	FROM sequences
	WHERE id = $1`

	var sequence Sequence
	var actions []SequenceAction

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&sequence.ID,
		&sequence.Name,
		&sequence.Description,
		&actions,
		&sequence.CreatedAt,
		&sequence.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	sequence.Actions = actions

	return &sequence, nil
}

func (m SequenceModel) Update(sequence *Sequence) error {
	query := `UPDATE sequences
	SET name = $2, description = $3, actions = $4, version = version + 1
	WHERE id = $1
	RETURNING version`

	args := []any{sequence.ID, sequence.Name, sequence.Description, sequence.Actions}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&sequence.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m SequenceModel) Delete(id uuid.UUID) error {
	query := `DELETE FROM sequences
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}
