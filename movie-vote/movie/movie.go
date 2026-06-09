package movie

import (
	"github.com/google/uuid"
)

type Movie struct {
	ID          string
	PollID      string
	Title       string
	ReleaseYear int
	Description string
}

type CreateMovieInput struct {
	Title       string
	PollID      string
	ReleaseYear int
	Description string
}

func CreateNewMovie(input CreateMovieInput) Movie {
	return Movie{
		ID:          uuid.New().String(),
		PollID:      input.PollID,
		Title:       input.Title,
		ReleaseYear: input.ReleaseYear,
		Description: input.Description,
	}
}
