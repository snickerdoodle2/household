package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"inzynierka/internal/data/validator"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorType string
type SensorInitBuffer map[uuid.UUID]Sensor

var (
	ErrDuplicateUri = errors.New("duplicate uri")
)

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

type SensorReturn interface {
	int | float64 | bool
}

type Sensor struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	URI         string     `json:"uri"`
	Type        SensorType `json:"type"`
	Hidden      bool       `json:"hidden"`
	RefreshRate int        `json:"refresh_rate"`
	CreatedAt   time.Time  `json:"created_at"`
	Version     int        `json:"version"`
	Active      bool       `json:"active"`
	IdToken     uuid.UUID  `json:"idToken"`
}

func ValidateSensor(v *validator.Validator, sensor *Sensor) {
	v.Check(sensor.Name != "", "name", "must be provided")
	v.Check(len(sensor.Name) <= 255, "name", "must not be more than 255 bytes long")

	v.Check(sensor.URI != "", "uri", "must be provided")
	v.Check(validator.Matches(sensor.URI, validator.UriRX), "uri", "must be valid uri")

	v.Check(sensor.Type != "", "type", "must be provided")
	v.Check(validator.PermittedValue(sensor.Type, SensorTypes...), "type", "must be known")

	isActive := sensor.Active
	v.Check(isActive || !isActive, "active", "must be a boolean")

	if !isActive {
		v.Check(sensor.RefreshRate != 0, "refresh_rate", "must be provided")
		v.Check(sensor.RefreshRate > 0, "refresh_rate", "must be a positive integer")
	}
}

type SensorModel struct {
	DB *pgxpool.Pool
}

func (m SensorModel) Insert(sensor *Sensor) error {
	query := `
    INSERT INTO sensors (id, name, uri, sensor_type, hidden, refresh_rate, active, id_token)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING created_at, version
    `

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	sensor.ID = uuid

	args := []any{sensor.ID, sensor.Name, sensor.URI, sensor.Type, sensor.Hidden, sensor.RefreshRate, sensor.Active, sensor.IdToken}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// TODO: handle not unique uuid edge case

	err = m.DB.QueryRow(ctx, query, args...).Scan(&sensor.CreatedAt, &sensor.Version)

	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint \"uri_unique\""):
			return ErrDuplicateUri
		default:
			return err
		}
	}

	return nil
}

func (m SensorModel) Get(id uuid.UUID) (*Sensor, error) {
	query := `
    SELECT id, name, uri, sensor_type, hidden, refresh_rate, created_at, version, active, id_token
    FROM sensors
    WHERE id = $1
    `
	var sensor Sensor

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&sensor.ID,
		&sensor.Name,
		&sensor.URI,
		&sensor.Type,
		&sensor.Hidden,
		&sensor.RefreshRate,
		&sensor.CreatedAt,
		&sensor.Version,
		&sensor.Active,
		&sensor.IdToken,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &sensor, nil
}

type SensorSimple struct {
	ID     uuid.UUID  `json:"id"`
	Name   string     `json:"name"`
	Type   SensorType `json:"type"`
	Hidden bool       `json:"hidden"`
	Active bool       `json:"active"`
}

func (m SensorModel) GetAllInfo() ([]*SensorSimple, error) {
	// TODO: add filtering and pagination
	query := `
    SELECT id, name, sensor_type, hidden, active
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

	sensors := []*SensorSimple{}

	for rows.Next() {
		var sensor SensorSimple

		err := rows.Scan(
			&sensor.ID,
			&sensor.Name,
			&sensor.Type,
			&sensor.Hidden,
			&sensor.Active,
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

func (m SensorModel) GetAll() ([]*Sensor, error) {
	query := `
    SELECT id, name, uri, sensor_type, hidden, refresh_rate, created_at, version, active, id_token
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
			&sensor.URI,
			&sensor.Type,
			&sensor.Hidden,
			&sensor.RefreshRate,
			&sensor.CreatedAt,
			&sensor.Version,
			&sensor.Active,
			&sensor.IdToken,
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

func (m SensorModel) GetUri(id uuid.UUID) (string, error) {
	query := `
    SELECT uri FROM sensors
    WHERE id = $1
    LIMIT 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uri string

	err := m.DB.QueryRow(ctx, query, id).Scan(&uri)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}
	return uri, nil
}

func (m SensorModel) Update(sensor *Sensor) error {
	query := `
    UPDATE sensors
    SET name = $1, uri = $2, sensor_type = $3, hidden = $4, refresh_rate = $5, active = $6, id_token = $7, version = version + 1
    WHERE id = $8
    RETURNING version
    `

	args := []any{
		sensor.Name,
		sensor.URI,
		sensor.Type,
		sensor.Hidden,
		sensor.RefreshRate,
		sensor.Active,
		sensor.IdToken,
		sensor.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&sensor.Version)

	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint \"uri_unique\""):
			return ErrDuplicateUri
		default:
			return err
		}
	}

	return nil
}

func (m SensorModel) DeleteSensorAndMeasurements(id uuid.UUID) error {
	deleteMeasurementsQuery := `
        DELETE FROM sensor_measurements
		WHERE sensor_id = $1;
    `

	deleteSensorQuery := `
		DELETE FROM sensors
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, deleteMeasurementsQuery, id)
	if err != nil {
		return err
	}

	result, err := m.DB.Exec(ctx, deleteSensorQuery, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m SensorModel) GetByIdToken(idToken uuid.UUID) (*Sensor, error) {
	query := `
    SELECT id, name, uri, sensor_type, hidden, refresh_rate, created_at, version, active, id_token
	FROM sensors
	WHERE id_token = $1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sensor Sensor
	err := m.DB.QueryRow(ctx, query, idToken).Scan(
		&sensor.ID,
		&sensor.Name,
		&sensor.URI,
		&sensor.Type,
		&sensor.Hidden,
		&sensor.RefreshRate,
		&sensor.CreatedAt,
		&sensor.Version,
		&sensor.Active,
		&sensor.IdToken,
	)

	if err == sql.ErrNoRows {
		return &sensor, fmt.Errorf("no sensor found with id token %s", idToken)
	}

	if err != nil {
		return &sensor, err
	}

	return &sensor, nil
}
