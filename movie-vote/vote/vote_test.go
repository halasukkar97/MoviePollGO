package vote

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestCreateNewVote(t *testing.T) {
	input := CreateVoteInput{
		PollID:   "poll-123",
		UserID:   "user-456",
		MovieIDs: []string{"movie-1", "movie-2"},
	}

	v := CreateNewVote(input)

	if _, err := uuid.Parse(v.ID); err != nil {
		t.Fatalf("expected valid UUID, got %q: %v", v.ID, err)
	}

	if v.PollID != input.PollID {
		t.Errorf("expected poll ID %q, got %q", input.PollID, v.PollID)
	}

	if v.UserID != input.UserID {
		t.Errorf("expected user ID %q, got %q", input.UserID, v.UserID)
	}

	if !reflect.DeepEqual(v.MovieIDs, input.MovieIDs) {
		t.Errorf("expected movie IDs %v, got %v", input.MovieIDs, v.MovieIDs)
	}
}
