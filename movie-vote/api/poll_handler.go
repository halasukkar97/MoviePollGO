package api

import (
	"encoding/json"
	"movie-vote/poll"
	"net/http"
	"time"
)

// CreatePollRequest Request body for creating a poll
type CreatePollRequest struct {
	Name              string    `json:"name"`
	MaxVotesPerPerson int       `json:"maxVotesPerPerson"`
	Deadline          time.Time `json:"deadline"`
}

// CreatePollResponse Response returned after creating a poll
type CreatePollResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	MaxVotesPerPerson int       `json:"maxVotesPerPerson"`
	IsClosed          bool      `json:"isClosed"`
	Deadline          time.Time `json:"deadline"`
}

// CreatePollHandler Handle POST /polls
func CreatePollHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest

	// Read JSON request body
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	// Create poll object
	createdPoll := poll.CreateNewPoll(poll.CreatePollInput{
		Name:              req.Name,
		MaxVotesPerPerson: req.MaxVotesPerPerson,
		Deadline:          req.Deadline,
	})

	// Save poll in memory
	polls = append(polls, createdPoll)

	// Build response
	response := CreatePollResponse{
		ID:                createdPoll.ID,
		Name:              createdPoll.Name,
		MaxVotesPerPerson: createdPoll.MaxVotesPerPerson,
		IsClosed:          createdPoll.IsClosed,
		Deadline:          createdPoll.Deadline,
	}

	// Return 201 Created
	w.WriteHeader(http.StatusCreated)

	// Send JSON response
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

// PollsHandler Route requests to the correct handler
func PollsHandler(w http.ResponseWriter, r *http.Request) {
	// Create poll
	if r.Method == http.MethodPost {
		CreatePollHandler(w, r)
		return
	}

	// List polls
	if r.Method == http.MethodGet {
		ListPollsHandler(w, r)
		return
	}

	// Reject unsupported methods
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// ListPollsHandler Handle GET /polls
func ListPollsHandler(w http.ResponseWriter, r *http.Request) {
	// Return all polls as JSON
	err := json.NewEncoder(w).Encode(polls)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	pollID := r.URL.Query().Get("pollId")

	foundPoll, found := FindPollByID(pollID)

	if !found {
		http.Error(w, "poll not found", http.StatusNotFound)
		return
	}

	results := foundPoll.GetResults()

	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
