package api

import (
	"encoding/json"
	"movie-vote/poll"
	"net/http"
	"time"
)

// CreatePollRequest is the JSON body clients send when they create a poll.
type CreatePollRequest struct {
	Name              string    `json:"name"`
	MaxVotesPerPerson int       `json:"maxVotesPerPerson"`
	Deadline          time.Time `json:"deadline"`
}

// CreatePollResponse is the JSON response sent back after a poll is created.
type CreatePollResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	MaxVotesPerPerson int       `json:"maxVotesPerPerson"`
	IsClosed          bool      `json:"isClosed"`
	Deadline          time.Time `json:"deadline"`
}

// CreatePollHandler handles POST /polls.
// It reads JSON from the request, creates a poll model, stores it, and returns it.
func CreatePollHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest

	// Decode turns the incoming JSON request body into a Go struct.
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	// The poll package owns the rules for building a new poll.
	createdPoll := poll.CreateNewPoll(poll.CreatePollInput{
		Name:              req.Name,
		MaxVotesPerPerson: req.MaxVotesPerPerson,
		Deadline:          req.Deadline,
	})

	// Save the poll so later requests can list it or add movies/votes to it.
	SavePoll(createdPoll)

	// Only expose the fields the API should return to the client.
	response := CreatePollResponse{
		ID:                createdPoll.ID,
		Name:              createdPoll.Name,
		MaxVotesPerPerson: createdPoll.MaxVotesPerPerson,
		IsClosed:          createdPoll.IsClosed,
		Deadline:          createdPoll.Deadline,
	}

	// StatusCreated means the request succeeded and created a new resource.
	w.WriteHeader(http.StatusCreated)

	// Encode writes the Go response struct back to the client as JSON.
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

// PollsHandler routes /polls requests by HTTP method.
func PollsHandler(w http.ResponseWriter, r *http.Request) {
	// POST /polls creates a new poll.
	if r.Method == http.MethodPost {
		CreatePollHandler(w, r)
		return
	}

	// GET /polls lists the polls currently in memory.
	if r.Method == http.MethodGet {
		ListPollsHandler(w, r)
		return
	}

	// Any other method, such as PUT or DELETE, is not supported for this route.
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// ListPollsHandler handles GET /polls.
func ListPollsHandler(w http.ResponseWriter, r *http.Request) {
	// Return all polls as JSON.
	err := json.NewEncoder(w).Encode(GetAllPolls())
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ResultsHandler handles GET /results?pollId=...
// It finds the requested poll and returns vote totals keyed by movie ID.
func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Query parameters come from the URL after the question mark.
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
