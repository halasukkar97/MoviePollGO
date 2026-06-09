package api

import (
	"encoding/json"
	"movie-vote/vote"
	"net/http"
)

type createVoteRequest struct {
	MovieIDs []string `json:"movieIds"`
	PollID   string   `json:"pollId"`
	UserID   string   `json:"userId"`
}

type createVoteResponse struct {
	ID       string   `json:"id"`
	PollID   string   `json:"pollId"`
	UserID   string   `json:"userId"`
	MovieIDs []string `json:"movieIds"`
}

func CreateVoteHandler(w http.ResponseWriter, r *http.Request) {
	var req createVoteRequest

	// Read request
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	createdVote := vote.CreateNewVote(vote.CreateVoteInput{
		PollID:   req.PollID,
		UserID:   req.UserID,
		MovieIDs: req.MovieIDs,
	})

	foundPoll, found := FindPollByID(req.PollID)

	if !found {
		http.Error(w, "poll not found", http.StatusNotFound)
		return
	}

	submitErr := foundPoll.SubmitVote(createdVote)
	if submitErr != nil {
		http.Error(w, submitErr.Error(), http.StatusBadRequest)
		return
	}

	votes = append(votes, createdVote)

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
