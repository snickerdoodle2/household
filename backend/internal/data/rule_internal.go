package data

import (
	"errors"
	"inzynierka/internal/data/validator"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
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

	ErrParseInvalidData = errors.New(`Field contains invalid data`)
)

type RuleData map[uuid.UUID]float64

type RuleInternal interface {
	Process(data RuleData, model *SensorMeasurementModel) (bool, error)
	// NOTE: map[uuid.UUID]struct{} (hashset) -> better perf
	// should be called only once per rule lifetime
	Dependencies() []uuid.UUID
	Validate(v *validator.Validator)
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

func unmarshalPerc(data map[string]interface{}) (*RulePerc, error) {
	idData, ok := data["sensor_id"]
	if !ok {
		return nil, ErrParseMissingSensorID
	}
	idStr, ok := idData.(string)
	if !ok {
		log.Error("idStr")
		return nil, ErrParseInvalidType
	}
	sensorID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	dur, ok := data["duration"]
	if !ok {
		return nil, ErrMissingVal
	}
	durStr, ok := dur.(string)
	if !ok {
		log.Error("durStr")
		return nil, ErrParseInvalidType
	}
	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return nil, err
	}

	percData, ok := data["perc"]
	if !ok {
		return nil, ErrMissingVal
	}
	perc, ok := percData.(float64)
	if !ok {
		return nil, ErrParseInvalidType
	}

	return &RulePerc{SensorID: sensorID, Delta: Duration(duration), Percentile: int(perc)}, nil
}

func unmarhsalField[T any](fieldName string, data map[string]interface{}) (*T, error) {
	str, ok := data[fieldName]
	if !ok {
		log.Error("unmarshalField", "field", fieldName, "action", "map")
		return nil, ErrParseMissingValue
	}
	value, ok := str.(T)
	if !ok {
		return nil, ErrParseInvalidType
	}

	return &value, nil
}

func unmarshalTime(data map[string]interface{}) (*RuleTime, error) {
	hour, err := unmarhsalField[float64]("hour", data)
	if err != nil {
		return nil, err
	}
	minute, err := unmarhsalField[float64]("minute", data)
	if err != nil {
		return nil, err
	}
	variant, err := unmarhsalField[string]("variant", data)
	if err != nil {
		return nil, err
	}

	return &RuleTime{Hour: int(*hour), Minute: int(*minute), Variant: TimeType(*variant)}, nil
}

func parseRuleDayField(field string, min, max int) ([]int, error) {
	if field == "*" {
		res := make([]int, max-min+1)
		for i := min; i <= max; i++ {
			res[i-min] = i
		}
		return res, nil
	}

	if strings.Contains(field, "-") {
		tmp := strings.Split(field, "-")

		left, err := strconv.Atoi(tmp[0])
		if err != nil {
			return nil, ErrParseInvalidData
		}

		right, err := strconv.Atoi(tmp[1])
		if err != nil {
			return nil, ErrParseInvalidData
		}

		if left > right {
			return nil, ErrParseInvalidData
		}

		res := make([]int, right-left+1)
		for i := left; i <= right; i++ {
			res[i-left] = i
		}
		return res, nil
	}

	x, err := strconv.Atoi(field)
	if err != nil {
		return nil, ErrParseInvalidData
	}

	return []int{x}, nil
}

func ParseRuleDay(format string) (*RuleDay, error) {
	fields := strings.Fields(format)
	if len(fields) != 3 {
		return nil, ErrParseInvalidData
	}

	days, err := parseRuleDayField(fields[0], 1, 31)
	if err != nil {
		return nil, err
	}

	monthsInt, err := parseRuleDayField(fields[1], 1, 12)
	if err != nil {
		return nil, err
	}

	months := make([]time.Month, len(monthsInt))
	for i, v := range monthsInt {
		months[i] = time.Month(v)
	}

	weekdaysInt, err := parseRuleDayField(fields[2], 1, 7)
	if err != nil {
		return nil, err
	}

	weekdays := make([]time.Weekday, len(weekdaysInt))
	for i, v := range weekdaysInt {
		if v == 7 {
			weekdays[i] = time.Sunday
		} else {
			weekdays[i] = time.Weekday(v)
		}
	}

	return &RuleDay{
		Format:   format,
		Days:     days,
		Months:   months,
		Weekdays: weekdays,
	}, nil
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

		child, err := UnmarshalInternalRuleJSON(wrapped)
		if err != nil {
			return nil, err
		}

		return &RuleNot{Wrapped: child}, nil
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
	case "perc":
		return unmarshalPerc(data)
	case "time":
		return unmarshalTime(data)
	default:
		return nil, ErrParseUnknownType
	}
}
