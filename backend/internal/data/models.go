package data

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	logger            = log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		Level:           log.InfoLevel,
		ReportCaller:    true,
	})
)

type Models struct {
	Sensors            SensorModel
	Rules              RuleModel
	Users              UserModel
	Tokens             TokenModel
	SensorMeasurements SensorMeasurementModel
	Sequences          SequenceModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Sensors:            SensorModel{DB: db},
		Rules:              RuleModel{DB: db},
		Users:              UserModel{DB: db},
		Tokens:             TokenModel{DB: db},
		SensorMeasurements: SensorMeasurementModel{DB: db},
		Sequences:          SequenceModel{DB: db},
	}
}
