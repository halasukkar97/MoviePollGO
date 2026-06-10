package api

import (
	"encoding/json"
	"movie-vote/movie"
	"net/http"
)

// CreateMovieRequest is the JSON body clients send when they add a movie to a poll.
type CreateMovieRequest struct {
	Title       string `json:"title"`
	PollID      string `json:"pollId"`
	ReleaseYear int    `json:"releaseYear"`
	Description string `json:"description"`
}

// CreateMovieResponse is the JSON response sent back after a movie is created.
type CreateMovieResponse struct {
	ID          string `json:"id"`
	PollID      string `json:"pollId"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseYear"`
	Description string `json:"description"`
}

// CreateMovieHandler handles POST /movies.
// It creates a movie, makes sure the target poll exists, and attaches the movie to that poll.
func CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateMovieRequest

	// Decode the JSON request body into a Go struct.
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	// Build the movie model before saving it.
	createdMovie := movie.CreateNewMovie(movie.CreateMovieInput{
		Title:       req.Title,
		PollID:      req.PollID,
		ReleaseYear: req.ReleaseYear,
		Description: req.Description,
	})

	// A movie must belong to an existing poll.
	foundPoll, found := FindPollByID(req.PollID)

	if !found {
		http.Error(w, "poll not found", http.StatusNotFound)
		return
	}

	// Keep a global list of movies for GET /movies.
	SaveMovie(createdMovie)

	// Also add the movie to its poll so voting validation can find it.
	foundPoll.AddMovie(createdMovie)

	// Return the created movie data to the client.
	response := CreateMovieResponse{
		ID:          createdMovie.ID,
		PollID:      createdMovie.PollID,
		Title:       createdMovie.Title,
		ReleaseYear: createdMovie.ReleaseYear,
		Description: createdMovie.Description,
	}

	// Send 201 Created before writing the JSON body.
	w.WriteHeader(http.StatusCreated)

	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

// MoviesHandler routes /movies requests by HTTP method.
func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	// POST /movies creates a new movie.
	if r.Method == http.MethodPost {
		CreateMovieHandler(w, r)
		return
	}

	// GET /movies lists every movie currently in memory.
	if r.Method == http.MethodGet {
		ListMoviesHandler(w, r)
		return
	}

	// Any other method is not supported for this route.
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// ListMoviesHandler handles GET /movies.
func ListMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Return all movies as JSON.
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
