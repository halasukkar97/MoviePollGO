package poll

import (
	"movie-vote/movie"
	"movie-vote/vote"
	"testing"
	"time"
)

func TestSubmitVoteSuccess(t *testing.T) {
	// Arrange
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	movie2 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)

	p.AddMovie(movie1)
	p.AddMovie(movie2)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
		},
	)

	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err != nil {
		t.Errorf("expected vote to succeed, got error: %v", err)
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
	// Arrange
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(-24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	movie2 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)

	p.AddMovie(movie1)
	p.AddMovie(movie2)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
		},
	)

	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err == nil {
		t.Errorf("expected poll expired error, got nil")
		return
	}

	if err.Error() != "poll has expired" {
		t.Errorf("expected poll has expired error, got: %v", err)
	}
}

func TestSubmitVotePollClosed(t *testing.T) {
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	movie2 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)

	p.IsClosed = true
	p.AddMovie(movie1)
	p.AddMovie(movie2)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
		},
	)

	p.IsClosed = true
	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err == nil {
		t.Errorf("expected poll expired error, got nil")
		return
	}

	if err.Error() != "poll is closed" {
		t.Errorf("expected poll has expired error, got: %v", err)
	}
}

func TestSubmitVoteAlreadyVoted(t *testing.T) {
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	movie2 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)

	movie3 := movie.CreateNewMovie(
		"titanic",
		p.ID,
		2021,
		"sea",
	)

	p.AddMovie(movie1)
	p.AddMovie(movie2)
	p.AddMovie(movie3)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
			movie3.ID,
		},
	)

	v2 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
			movie3.ID,
		},
	)
	p.SubmitVote(v2)

	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err == nil {
		t.Errorf("expected poll expired error, got nil")
		return
	}

	if err.Error() != "you have already voted for this poll" {
		t.Errorf("expected poll has expired error, got: %v", err)
	}
}

func TestSubmitVoteMovieDoesNotExist(t *testing.T) {
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	movie2 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)

	p.AddMovie(movie1)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
		},
	)

	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err == nil {
		t.Errorf("expected poll expired error, got nil")
		return
	}

	if err.Error() != "this movie doesn't exist in this poll" {
		t.Errorf("expected poll has expired error, got: %v", err)
	}
}

func TestSubmitVoteDuplicateMovie(t *testing.T) {
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	p.AddMovie(movie1)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie1.ID,
		},
	)

	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err == nil {
		t.Errorf("expected duplicate movie error, got nil), got nil")
		return
	}

	if err.Error() != "duplicated votes for the same movie are not allowed" {
		t.Errorf("expected poll has expired error, got: %v", err)
	}
}

func TestSubmitVoteTooManyMovies(t *testing.T) {
	p := CreateNewPoll(
		"Friday Movie Night",
		3,
		time.Now().Add(24*time.Hour),
	)

	movie1 := movie.CreateNewMovie(
		"Interstellar",
		p.ID,
		2014,
		"Space exploration",
	)

	movie2 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)
	movie3 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)
	movie4 := movie.CreateNewMovie(
		"Dune",
		p.ID,
		2021,
		"Arrakis",
	)

	p.AddMovie(movie1)
	p.AddMovie(movie2)
	p.AddMovie(movie3)
	p.AddMovie(movie4)

	v1 := vote.CreateNewVote(
		p.ID,
		"hela-user",
		[]string{
			movie1.ID,
			movie2.ID,
			movie3.ID,
			movie4.ID,
		},
	)

	// Act
	err := p.SubmitVote(v1)

	// Assert
	if err == nil {
		t.Errorf("expected poll expired error, got nil")
		return
	}

	if err.Error() != "too many movies selected" {
		t.Errorf("expected poll has expired error, got: %v", err)
	}
}
