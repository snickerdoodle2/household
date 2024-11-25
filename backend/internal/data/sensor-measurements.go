package data

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorMeasurement struct {
	SensorID      uuid.UUID
	MeasuredAt    time.Time
	MeasuredValue float64
}

type SensorMeasurementModel struct {
	DB *pgxpool.Pool
}

func (m *SensorMeasurementModel) Insert(measurement *SensorMeasurement) error {
	query := `
    INSERT INTO sensor_measurements (sensor_id, measured_at, measured_value)
    VALUES ($1, $2, $3)
    `

	args := []any{measurement.SensorID, measurement.MeasuredAt, measurement.MeasuredValue}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)

	if err != nil {
		return err
	}
	return nil
}

func (m *SensorMeasurementModel) GetLastMeasurement(id uuid.UUID) (*SensorMeasurement, error) {
	query := `
	SELECT measured_at, measured_value
	FROM sensor_measurements
	WHERE sensor_id = $1
	ORDER BY measured_at DESC
	LIMIT 1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lastMeasurement := SensorMeasurement{SensorID: id}
	err := m.DB.QueryRow(ctx, query, id).Scan(&lastMeasurement.MeasuredAt, &lastMeasurement.MeasuredValue)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &lastMeasurement, nil
}

func (m *SensorMeasurementModel) GetLastNMeasurements(id uuid.UUID, n int) ([]*SensorMeasurement, error) {
	query := `
	SELECT measured_at, measured_value
	FROM sensor_measurements
	WHERE sensor_id = $1
	ORDER BY measured_at DESC
	LIMIT $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{id, n}

	rows, err := m.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	measurements := []*SensorMeasurement{}

	for rows.Next() {
		measurement := SensorMeasurement{SensorID: id}

		err := rows.Scan(&measurement.MeasuredAt, &measurement.MeasuredValue)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, &measurement)
	}

	return measurements, nil
}

// INFO: how to name this :(
func (m *SensorMeasurementModel) GetMeasurementsSince(id uuid.UUID, delta time.Duration) ([]*SensorMeasurement, error) {
	query := `
    SELECT measured_at, measured_value from sensor_measurements
    WHERE sensor_id = $1
    AND now() - measured_at < $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, id, delta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	measurements := []*SensorMeasurement{}

	for rows.Next() {
		measurement := SensorMeasurement{SensorID: id}

		err := rows.Scan(&measurement.MeasuredAt, &measurement.MeasuredValue)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, &measurement)
	}

	return measurements, nil
}
