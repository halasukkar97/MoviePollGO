package poll

import (
	"movie-vote/movie"
	"movie-vote/vote"
	"testing"
	"time"
)

func newTestPoll(maxVotes int, deadline time.Time) Poll {
	return CreateNewPoll(CreatePollInput{
		Name:              "Friday Movie Night",
		MaxVotesPerPerson: maxVotes,
		Deadline:          deadline,
	})
}

func newTestMovie(pollID, title string) movie.Movie {
	return movie.CreateNewMovie(movie.CreateMovieInput{
		Title:       title,
		PollID:      pollID,
		ReleaseYear: 2021,
		Description: title + " description",
	})
}

func newTestVote(pollID, userID string, movieIDs []string) vote.Vote {
	return vote.CreateNewVote(vote.CreateVoteInput{
		PollID:   pollID,
		UserID:   userID,
		MovieIDs: movieIDs,
	})
}

func requireError(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error %q, got nil", expected)
	}

	if err.Error() != expected {
		t.Fatalf("expected error %q, got %q", expected, err.Error())
	}
}

func TestSubmitVoteSuccess(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")
	movie2 := newTestMovie(p.ID, "Dune")

	p.AddMovie(movie1)
	p.AddMovie(movie2)

	v := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie2.ID})

	err := p.SubmitVote(v)
	if err != nil {
		t.Fatalf("expected vote to succeed, got error: %v", err)
	}

	results := p.GetResults()
	if results[movie1.ID] != 1 {
		t.Errorf("expected Interstellar to have 1 vote, got %d", results[movie1.ID])
	}

	if results[movie2.ID] != 1 {
		t.Errorf("expected Dune to have 1 vote, got %d", results[movie2.ID])
	}
}

func TestSubmitVotePollExpired(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(-24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")
	movie2 := newTestMovie(p.ID, "Dune")

	p.AddMovie(movie1)
	p.AddMovie(movie2)

	v := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie2.ID})

	err := p.SubmitVote(v)

	requireError(t, err, "poll has expired")
}

func TestSubmitVotePollClosed(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")
	movie2 := newTestMovie(p.ID, "Dune")

	p.Close()
	p.AddMovie(movie1)
	p.AddMovie(movie2)

	v := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie2.ID})

	err := p.SubmitVote(v)

	requireError(t, err, "poll is closed")
}

func TestSubmitVoteAlreadyVoted(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")
	movie2 := newTestMovie(p.ID, "Dune")
	movie3 := newTestMovie(p.ID, "Titanic")

	p.AddMovie(movie1)
	p.AddMovie(movie2)
	p.AddMovie(movie3)

	firstVote := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie2.ID, movie3.ID})
	secondVote := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie2.ID, movie3.ID})

	if err := p.SubmitVote(firstVote); err != nil {
		t.Fatalf("expected first vote to succeed, got error: %v", err)
	}

	err := p.SubmitVote(secondVote)

	requireError(t, err, "you have already voted for this poll")
}

func TestSubmitVoteMovieDoesNotExist(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")
	movie2 := newTestMovie(p.ID, "Dune")

	p.AddMovie(movie1)

	v := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie2.ID})

	err := p.SubmitVote(v)

	requireError(t, err, "this movie doesn't exist in this poll")
}

func TestSubmitVoteDuplicateMovie(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")

	p.AddMovie(movie1)

	v := newTestVote(p.ID, "hela-user", []string{movie1.ID, movie1.ID})

	err := p.SubmitVote(v)

	requireError(t, err, "duplicated votes for the same movie are not allowed")
}

func TestSubmitVoteTooManyMovies(t *testing.T) {
	p := newTestPoll(3, time.Now().Add(24*time.Hour))
	movie1 := newTestMovie(p.ID, "Interstellar")
	movie2 := newTestMovie(p.ID, "Dune")
	movie3 := newTestMovie(p.ID, "Titanic")
	movie4 := newTestMovie(p.ID, "Arrival")

	p.AddMovie(movie1)
	p.AddMovie(movie2)
	p.AddMovie(movie3)
	p.AddMovie(movie4)

	v := newTestVote(p.ID, "hela-user", []string{
		movie1.ID,
		movie2.ID,
		movie3.ID,
		movie4.ID,
	})

	err := p.SubmitVote(v)

	requireError(t, err, "too many movies selected")
}
