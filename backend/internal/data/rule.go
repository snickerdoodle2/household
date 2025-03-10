package data

import (
	"context"
	"encoding/json"
	"errors"
	"inzynierka/internal/data/validator"
	"reflect"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNonExistingTo         = errors.New("on_valid.id is not pointing to any existing device")
	ErrMissingDependencyChan = errors.New("Provided listeners map is missing required dependency")
)

type TargetType string

const (
	SensorTarget   TargetType = "sensor"
	SequenceTarget TargetType = "sequence"
)

type ValidRuleAction struct {
	TargetType TargetType             `json:"target_type"`
	TargetId   uuid.UUID              `json:"target_id"`
	Payload    map[string]interface{} `json:"payload"`
}

type Rule struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Internal    RuleInternal    `json:"internal"`
	OnValid     ValidRuleAction `json:"on_valid"`
	CreatedAt   time.Time       `json:"created_at"`
	Version     int             `json:"version"`
	prev        bool
}

type SensorListeners map[uuid.UUID]*Listener[float64]

func (t TargetType) IsValid() bool {
	return t == SensorTarget || t == SequenceTarget
}

// TOOD: Handle stopping on channel close
// REF: https://pkg.go.dev/reflect#Select
func (r *Rule) Run(listeners SensorListeners, validCh chan ValidRuleAction, stopCh chan struct{}, m *SensorMeasurementModel) error {
	deps := r.Internal.Dependencies()
	channels := make([]reflect.SelectCase, len(deps)+1)
	values := make(RuleData)
	for i, dep := range deps {
		listener, ok := listeners[dep]
		if !ok {
			return ErrMissingDependencyChan
		}
		cur := listener.GetCurrentValue()
		if len(cur) > 0 {
			values[dep] = cur[len(cur)-1]
		}

		msgCh := listener.GetBroker().Subscribe()
		defer listener.GetBroker().Unsubscribe(msgCh)

		channels[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(msgCh)}
	}
	channels[len(deps)] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(stopCh)}

	r.update(values, validCh, m)

	for {
		i, sliceV, ok := reflect.Select(channels)
		if !ok {
			break
		}

		if i == len(deps) { // STOP CHANNEL
			logger.Debug("stopping rule")
			break
		}
		slice := sliceV.Interface().([]float64)
		values[deps[i]] = slice[len(slice)-1]
		// updating rule, sending onValid struct to channel if the rule has just been fulfilled
		r.update(values, validCh, m)
	}
	return nil
}

func (r *Rule) update(data RuleData, ch chan ValidRuleAction, m *SensorMeasurementModel) {
	cur, err := r.Internal.Process(data, m)
	if err != nil {
		return
	}

	// If something changed from previous
	if cur != r.prev {
		if cur {
			ch <- r.OnValid
		} else {
		}

		r.prev = cur
	}
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
	v.Check(r.OnValid.TargetType.IsValid(), "on_valid.target-type", "must be either 'sensor' or 'sequence'")
}

type RuleModel struct {
	DB *pgxpool.Pool
}

func (m *RuleModel) Insert(rule *Rule) error {
	query := `
    INSERT INTO rules (id, name, description, internal, valid_target_type, valid_target_id, valid_target_payload)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING created_at, version
    `

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	rule.ID = uuid

	args := []any{uuid, rule.Name, rule.Description, rule.Internal, rule.OnValid.TargetType, rule.OnValid.TargetId, rule.OnValid.Payload}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = m.DB.QueryRow(ctx, query, args...).Scan(&rule.CreatedAt, &rule.Version)

	if err != nil {
		return err
	}
	return nil
}

func (m *RuleModel) Get(id uuid.UUID) (*Rule, error) {
	query := `
    SELECT id, name, description, internal, valid_target_type, valid_target_id, valid_target_payload, created_at, version
    FROM rules
    WHERE id = $1
    `

	var ruleS Rule
	var internalMap map[string]interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&ruleS.ID,
		&ruleS.Name,
		&ruleS.Description,
		&internalMap,
		&ruleS.OnValid.TargetType,
		&ruleS.OnValid.TargetId,
		&ruleS.OnValid.Payload,
		&ruleS.CreatedAt,
		&ruleS.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	internal, err := UnmarshalInternalRuleJSON(internalMap)
	if err != nil {
		return nil, err
	}

	ruleS.Internal = internal

	return &ruleS, nil
}

func (m *RuleModel) GetAll() ([]*Rule, error) {
	query := `
    SELECT id, name, description, internal, valid_target_type, valid_target_id, valid_target_payload, created_at, version
    FROM rules
    ORDER BY id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := []*Rule{}

	for rows.Next() {
		var ruleS Rule
		var internalMap map[string]interface{}

		err = rows.Scan(
			&ruleS.ID,
			&ruleS.Name,
			&ruleS.Description,
			&internalMap,
			&ruleS.OnValid.TargetType,
			&ruleS.OnValid.TargetId,
			&ruleS.OnValid.Payload,
			&ruleS.CreatedAt,
			&ruleS.Version,
		)

		if err != nil {
			return nil, err
		}

		internal, err := UnmarshalInternalRuleJSON(internalMap)
		if err != nil {
			return nil, err
		}

		ruleS.Internal = internal

		rules = append(rules, &ruleS)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

type RuleSimple struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (m RuleModel) GetAllInfo() ([]*RuleSimple, error) {
	query := `
    SELECT id, name, description
    FROM rules
    ORDER BY id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := []*RuleSimple{}

	for rows.Next() {
		var ruleS RuleSimple

		err = rows.Scan(
			&ruleS.ID,
			&ruleS.Name,
			&ruleS.Description,
		)

		if err != nil {
			return nil, err
		}
		rules = append(rules, &ruleS)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (m RuleModel) Update(rule *Rule) error {
	query := `
       UPDATE rules
       SET name = $1, description = $2, internal = $3, valid_target_type = $4, valid_target_id = $5, valid_target_payload = $6, version = version + 1
       WHERE id = $7
       RETURNING version 
    `

	args := []any{
		rule.Name,
		rule.Description,
		rule.Internal,
		rule.OnValid.TargetType,
		rule.OnValid.TargetId,
		rule.OnValid.Payload,
		rule.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&rule.Version)

	if err != nil {
		return err
	}
	return nil
}

func (m RuleModel) Delete(id uuid.UUID) error {
	query := `
        DELETE FROM rules
        WHERE id = $1
    `

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
