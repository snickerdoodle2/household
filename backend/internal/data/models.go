package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Sensors SensorModel
	Rules   RuleModel
	Users   UserModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Sensors: SensorModel{DB: db},
		Rules:   RuleModel{DB: db},
		Users:   UserModel{DB: db},
	}
}
