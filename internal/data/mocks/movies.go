package mocks

import (
	"time"

	"github.com/juliflorezg/greenlight/internal/data"
)

var mockMovie = data.Movie{
	ID:        1,
	Title:     "The Hunger Games",
	Runtime:   142,
	Year:      2012,
	Genres:    []string{"dystopian sci-fi", "action", "adventure"},
	CreatedAt: time.Now(),
	Version:   1,
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *data.Movie) error {
	movie.ID = 1
	movie.CreatedAt = time.Now()
	movie.Version = 1

	return nil
}

func (m MockMovieModel) Get(id int64) (*data.Movie, error) {
	// // todo: Mock the action...
	switch id {
	case 1:
		return &mockMovie, nil
	default:
		return nil, data.ErrRecordNotFound
	}

	// return nil, nil
}

func (m MockMovieModel) Update(movie *data.Movie) error {
	// todo: Mock the action...
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	// todo: Mock the action...
	return nil
}

func NewMockModels() data.Models {
	return data.Models{
		Movies: MockMovieModel{},
	}
}
