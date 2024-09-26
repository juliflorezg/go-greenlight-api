package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("unable to update resource due to edit conflict")
)

type MoviesInterface interface {
	Insert(movie *Movie) error
	Get(id int64) (*Movie, error)
	Update(movie *Movie) error
	Delete(id int64) error
}

type Models struct {
	Movies MoviesInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
