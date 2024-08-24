package rule

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ValidRuleAction struct {
	To      uuid.UUID         `json:"to"`
	Payload map[string]string `json:"payload"`
}

type Rule struct {
	ID          uuid.UUID       `json:"id"`
	Description string          `json:"description"`
	Internal    RuleInternal    `json:"internal"`
	OnValid     ValidRuleAction `json:"on_valid"`
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID          uuid.UUID              `json:"id"`
		Description string                 `json:"description"`
		Internal    map[string]interface{} `json:"internal"`
		OnValid     ValidRuleAction        `json:"on_valid"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	r.ID = tmp.ID
	r.Description = tmp.Description
	r.OnValid = tmp.OnValid

	internal, err := UnmarshalInternalRuleJSON(tmp.Internal)
	if err != nil {
		return err
	}

	r.Internal = internal

	return nil
}
