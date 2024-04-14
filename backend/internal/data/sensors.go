package data

import (
	"context"
	"inzynierka/internal/data/validator"
	"time"

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

func (m SensorModel) Insert(sensor *Sensor) error {
	query := `
    INSERT INTO sensors (id, name, uri, sensor_type, refresh_rate)
    VALUES ($1, $2, $3, $4, $5)
    `

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	sensor.ID = uuid

	args := []any{sensor.ID, sensor.Name, sensor.URI, sensor.Type, sensor.RefreshRate}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = m.DB.Exec(ctx, query, args...)
	return err
}

func (m SensorModel) GetAll() ([]*Sensor, error) {
	// TODO: add filtering and pagination
	query := `
    SELECT id, name, sensor_type
    FROM sensors
    ORDER BY id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sensors := []*Sensor{}

	for rows.Next() {
		var sensor Sensor

		err := rows.Scan(
			&sensor.ID,
			&sensor.Name,
			&sensor.Type,
		)

		if err != nil {
			return nil, err
		}

		sensors = append(sensors, &sensor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

func (m SensorModel) Update(Sensor *Sensor) error {
	return nil
}

func (m SensorModel) Delete(id uuid.UUID) error {
	return nil
}
