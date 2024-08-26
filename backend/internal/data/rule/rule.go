package rule

import (
	"context"
	"encoding/json"
	"errors"
	"inzynierka/internal/data/validator"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNonExistingTo = errors.New("on_valid.id is not pointing to any existing device")
)

type ValidRuleAction struct {
	To      uuid.UUID         `json:"to"`
	Payload map[string]string `json:"payload"`
}

type Rule struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Internal    RuleInternal    `json:"internal"`
	OnValid     ValidRuleAction `json:"on_valid"`
	CreatedAt   time.Time       `json:"created_at"`
	Version     int             `json:"version"`
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID          uuid.UUID              `json:"id"`
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Internal    map[string]interface{} `json:"internal"`
		OnValid     ValidRuleAction        `json:"on_valid"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	r.ID = tmp.ID
	r.Name = tmp.Name
	r.Description = tmp.Description
	r.OnValid = tmp.OnValid

	internal, err := UnmarshalInternalRuleJSON(tmp.Internal)
	if err != nil {
		return err
	}

	r.Internal = internal

	return nil
}

func ValidateRule(v *validator.Validator, r *Rule) {
	v.Check(utf8.RuneCountInString(r.Name) > 0, "name", "must not be empty")
	v.Check(utf8.RuneCountInString(r.Name) <= 32, "name", "must not be longer than 32 characters")
	v.Check(utf8.RuneCountInString(r.Description) <= 256, "description", "must not be longer than 256 characters")
}

type RuleModel struct {
	DB *pgxpool.Pool
}

func (m *RuleModel) Insert(rule *Rule) error {
	query := `
    INSERT INTO rules (id, name, description, internal, valid_sensor_id, valid_payload)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at, version
    `

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	rule.ID = uuid

	args := []any{uuid, rule.Name, rule.Description, rule.Internal, rule.OnValid.To, rule.OnValid.Payload}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = m.DB.QueryRow(ctx, query, args...).Scan(&rule.CreatedAt, &rule.Version)

	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: insert or update on table \"rules\" violates foreign key constraint \"rules_valid_sensor_id_fkey\""):
			return ErrNonExistingTo
		default:
			return err
		}
	}
	return nil
}
