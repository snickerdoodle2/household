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
	sensor_id      uuid.UUID
	measured_at    time.Time
	measured_value float64
}

type SensorMeasurementModel struct {
	DB *pgxpool.Pool
}

func (m *SensorMeasurementModel) Insert(measurement *SensorMeasurement) error {
	query := `
    INSERT INTO sensor_measurements (sensor_id, measured_at, measured_value,)
    VALUES ($1, $2, $3)
    `

	args := []any{measurement.sensor_id, measurement.measured_at, measurement.measured_value}

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

	lastMeasurement := SensorMeasurement{sensor_id: id}
	err := m.DB.QueryRow(ctx, query, id).Scan(&lastMeasurement.measured_at, &lastMeasurement.measured_value)

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
		measurement := SensorMeasurement{sensor_id: id}

		err := rows.Scan(&measurement.measured_at, &measurement.measured_value)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, &measurement)
	}

	return measurements, nil
}
