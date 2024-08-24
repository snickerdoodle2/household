package rule

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// TODO: Add stack trace / sensor id to this erorr value
	ErrMissingVal           = errors.New("Missing value in data")
	ErrParseMissingType     = errors.New(`Missing "type" field in provided data`)
	ErrParseMissingChildren = errors.New(`Missing "children" field in provided data`)

	ErrParseMissingSensorID = errors.New(`Missing "sensor_id" field in provided data`)
	ErrParseMissingValue    = errors.New(`Missing "value" field in provided data`)

	ErrParseMissingWrapped = errors.New(`Missing "wrapped" field in provided data`)

	ErrParseInvalidType = errors.New(`Field is invalid type`)
	ErrParseUnknownType = errors.New(`Unknown "type" field in provided data`)
)

type RuleData map[uuid.UUID]float64

type RuleInternal interface {
	Process(data RuleData) (bool, error)
	// NOTE: map[uuid.UUID]struct{} (hashset) -> better perf
	// should be called only once per rule lifetime
	Dependencies() []uuid.UUID
}

func unmarshalChildren(data map[string]interface{}) ([]RuleInternal, error) {
	childrenData, ok := data["children"]
	if !ok {
		return nil, ErrParseMissingChildren
	}
	childrenArr := childrenData.([]interface{})

	res := make([]RuleInternal, 0)

	for _, childT := range childrenArr {
		child, ok := childT.(map[string]interface{})
		if !ok {
			return nil, ErrParseInvalidType
		}
		childParsed, err := UnmarshalInternalRuleJSON(child)
		if err != nil {
			return nil, err
		}
		res = append(res, childParsed)
	}

	return res, nil
}

func unmarshalSimple(data map[string]interface{}) (uuid.UUID, float64, error) {
	idData, ok := data["sensor_id"]
	if !ok {
		return uuid.UUID{}, 0, ErrParseMissingSensorID
	}

	idStr, ok := idData.(string)
	if !ok {
		return uuid.UUID{}, 0, ErrParseInvalidType
	}

	sensorID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	valueData, ok := data["value"]
	if !ok {
		return uuid.UUID{}, 0, ErrParseMissingValue
	}

	value, ok := valueData.(float64)
	if !ok {
		return uuid.UUID{}, 0, ErrParseInvalidType
	}

	return sensorID, value, nil
}

func UnmarshalInternalRuleJSON(data map[string]interface{}) (RuleInternal, error) {
	nodeType, ok := data["type"]

	if !ok {
		return nil, ErrParseMissingType
	}

	switch nodeType {
	case "and":
		children, err := unmarshalChildren(data)
		if err != nil {
			return nil, err
		}

		return &RuleAnd{Children: children}, nil
	case "or":
		children, err := unmarshalChildren(data)
		if err != nil {
			return nil, err
		}

		return &RuleOr{Children: children}, nil
	case "not":
		wrappedData, ok := data["wrapped"]
		if !ok {
			return nil, ErrParseMissingWrapped
		}

		wrapped, ok := wrappedData.(map[string]interface{})
		if !ok {
			return nil, ErrParseInvalidType
		}

		return UnmarshalInternalRuleJSON(wrapped)
	case "gt":
		sensorID, value, err := unmarshalSimple(data)
		if err != nil {
			return nil, err
		}

		return &RuleGT{SensorID: sensorID, Value: value}, nil
	case "lt":
		sensorID, value, err := unmarshalSimple(data)
		if err != nil {
			return nil, err
		}

		return &RuleLT{SensorID: sensorID, Value: value}, nil
	default:
		return nil, ErrParseUnknownType
	}
}
