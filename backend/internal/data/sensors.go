package data

import (
	"inzynierka/internal/data/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorType string

const (
	BinarySwitch  SensorType = "binary_switch"
	BinarySensor  SensorType = "binary_sensor"
	DecimalSwitch SensorType = "decimal_switch"
	DecimalSensor SensorType = "decimal_sensor"
	Button        SensorType = "button"
)

var SensorTypes = []SensorType{
	BinarySwitch,
	BinarySensor,
	DecimalSwitch,
	DecimalSensor,
	Button,
}

type Sensor struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	URI         string     `json:"-"`
	Type        SensorType `json:"type"`
	RefreshRate int        `json:"-"`
}

func ValidateSensor(v *validator.Validator, sensor *Sensor) {
	v.Check(sensor.Name != "", "name", "must be provided")
	v.Check(len(sensor.Name) <= 255, "name", "must not be more than 255 bytes long")

	v.Check(sensor.URI != "", "uri", "must be provided")
	v.Check(validator.Matches(sensor.URI, validator.UriRX), "uri", "must be valid uri")

	v.Check(sensor.Type != "", "type", "must be provided")
	v.Check(validator.PermittedValue(sensor.Type, SensorTypes...), "type", "must be known")

	v.Check(sensor.RefreshRate != 0, "refresh_rate", "must be provided")
	v.Check(sensor.RefreshRate > 0, "refresh_rate", "must be a positive integer")
}

type SensorModel struct {
	DB *pgxpool.Pool
}
