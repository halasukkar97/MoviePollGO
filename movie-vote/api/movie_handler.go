package api

import (
	"encoding/json"
	"movie-vote/movie"
	"net/http"
)

// CreateMovieRequest Request body for creating a movie
type CreateMovieRequest struct {
	Title       string `json:"title"`
	PollID      string `json:"pollId"`
	ReleaseYear int    `json:"releaseYear"`
	Description string `json:"description"`
}

// CreateMovieResponse Response returned after creating a movie
type CreateMovieResponse struct {
	ID          string `json:"id"`
	PollID      string `json:"pollId"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
	Description string `json:"description"`
}

// CreateMovieHandler handles POST /movies.
func CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateMovieRequest

	// Read request
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	// Create movie
	createdMovie := movie.CreateNewMovie(movie.CreateMovieInput{
		Title:       req.Title,
		PollID:      req.PollID,
		ReleaseYear: req.ReleaseYear,
		Description: req.Description,
	})

	// Find poll
	foundPoll, found := FindPollByID(req.PollID)

	if !found {
		http.Error(w, "poll not found", http.StatusNotFound)
		return
	}

	// Save movie
	movies = append(movies, createdMovie)

	// Add movie to poll
	foundPoll.AddMovie(createdMovie)

	// Build response
	response := CreateMovieResponse{
		ID:          createdMovie.ID,
		PollID:      createdMovie.PollID,
		Title:       createdMovie.Title,
		ReleaseYear: createdMovie.ReleaseYear,
		Description: createdMovie.Description,
	}

	// Send response
	w.WriteHeader(http.StatusCreated)

	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

// MoviesHandler Route requests to the correct handler
func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Create movie
	if r.Method == http.MethodPost {
		CreateMovieHandler(w, r)
		return
	}

	// List movies
	if r.Method == http.MethodGet {
		ListMoviesHandler(w, r)
		return
	}

	// Reject unsupported methods
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// ListMoviesHandler Handle GET /movies
func ListMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Return all movies as JSON
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
