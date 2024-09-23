package mocks

import (
	"github.com/juliflorezg/greenlight/internal/data"
)

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *data.Movie) error {
	// todo: Mock the action...
	return nil
}

func (m MockMovieModel) Get(id int64) (*data.Movie, error) {
	// todo: Mock the action...
	return nil, nil
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
