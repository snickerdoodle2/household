package data

import (
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

type Sensor struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	URI         string     `json:"-"`
	Type        SensorType `json:"type"`
	RefreshRate int        `json:"-"`
}

type SensorModel struct {
	DB *pgxpool.Pool
}
