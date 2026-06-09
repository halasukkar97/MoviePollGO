package vote

import (
	"github.com/google/uuid"
)

type Vote struct {
	ID       string
	PollID   string
	UserID   string
	MovieIDs []string
}

type CreateVoteInput struct {
	PollID   string
	UserID   string
	MovieIDs []string
}

func CreateNewVote(input CreateVoteInput) Vote {
	return Vote{
		ID:       uuid.New().String(),
		PollID:   input.PollID,
		UserID:   input.UserID,
		MovieIDs: input.MovieIDs,
	}
}
