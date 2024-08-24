package rule

import (
	"encoding/json"
	"inzynierka/internal/data/validator"
	"unicode/utf8"

	"github.com/google/uuid"
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
