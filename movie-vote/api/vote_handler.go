package api

import (
	"encoding/json"
	"movie-vote/vote"
	"net/http"
)

// createVoteRequest is the JSON body clients send when they vote in a poll.
type createVoteRequest struct {
	MovieIDs []string `json:"movieIds"`
	PollID   string   `json:"pollId"`
	UserID   string   `json:"userId"`
}

// createVoteResponse is the JSON response sent back after a vote is accepted.
type createVoteResponse struct {
	ID       string   `json:"id"`
	PollID   string   `json:"pollId"`
	UserID   string   `json:"userId"`
	MovieIDs []string `json:"movieIds"`
}

// CreateVoteHandler handles POST /votes.
// It creates a vote, asks the poll to validate it, stores it, and returns it.
func CreateVoteHandler(w http.ResponseWriter, r *http.Request) {
	var req createVoteRequest

	// Decode the JSON request body into a Go struct.
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	// Build the vote model from the request data.
	createdVote := vote.CreateNewVote(vote.CreateVoteInput{
		PollID:   req.PollID,
		UserID:   req.UserID,
		MovieIDs: req.MovieIDs,
	})

	// A vote can only be submitted to an existing poll.
	foundPoll, found := FindPollByID(req.PollID)

	if !found {
		http.Error(w, "poll not found", http.StatusNotFound)
		return
	}

	// SubmitVote checks the poll rules before adding the vote to the poll.
	submitErr := foundPoll.SubmitVote(createdVote)
	if submitErr != nil {
		http.Error(w, submitErr.Error(), http.StatusBadRequest)
		return
	}

	// Keep a global vote list as well as the vote stored inside the poll.
	SaveVote(createdVote)

	response := createVoteResponse{
		ID:       createdVote.ID,
		PollID:   createdVote.PollID,
		UserID:   createdVote.UserID,
		MovieIDs: createdVote.MovieIDs,
	}

	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}
}
