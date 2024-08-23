package rule

import "github.com/google/uuid"

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
