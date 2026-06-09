package poll

import (
	"errors"
	"movie-vote/movie"
	"movie-vote/vote"
	"time"

	"github.com/google/uuid"
)

// Poll represents a movie voting poll.
type Poll struct {
	ID                string
	Name              string
	IsClosed          bool
	MaxVotesPerPerson int
	Deadline          time.Time
	Movies            []movie.Movie
	Votes             []vote.Vote
}

type CreatePollInput struct {
	Name              string
	MaxVotesPerPerson int
	Deadline          time.Time
}

// CreateNewPoll creates a new poll.
func CreateNewPoll(input CreatePollInput) Poll {
	return Poll{
		ID:                uuid.New().String(),
		Name:              input.Name,
		MaxVotesPerPerson: input.MaxVotesPerPerson,
		Deadline:          input.Deadline,
		Movies:            []movie.Movie{},
		Votes:             []vote.Vote{},
	}
}

// AddMovie adds a movie to the poll.
func (p *Poll) AddMovie(movie movie.Movie) {
	p.Movies = append(p.Movies, movie)
}

// AddVote adds a vote to the poll.
func (p *Poll) AddVote(vote vote.Vote) {
	p.Votes = append(p.Votes, vote)
}

// Close marks the poll as closed.
func (p *Poll) Close() {
	p.IsClosed = true
}

// ValidateVoteCount checks if the vote count is allowed.
func (p *Poll) ValidateVoteCount(selectedMovieIDs []string) bool {
	return len(selectedMovieIDs) <= p.MaxVotesPerPerson
}

// GetResults returns the vote count per movie.
func (p *Poll) GetResults() map[string]int {
	result := make(map[string]int)

	for _, v := range p.Votes {
		for _, movieID := range v.MovieIDs {
			result[movieID]++
		}
	}

	return result
}

// HasMovie checks if a movie belongs to the poll.
func (p *Poll) HasMovie(movieID string) bool {
	for _, film := range p.Movies {
		if film.ID == movieID {
			return true
		}
	}

	return false
}

// HasDuplicateMovies checks for duplicate movie votes.
func (p *Poll) HasDuplicateMovies(v vote.Vote) bool {
	seenMovies := make(map[string]bool)

	for _, movieID := range v.MovieIDs {
		if seenMovies[movieID] {
			return true
		}

		seenMovies[movieID] = true
	}

	return false
}

// AlreadyVoted checks if a user has already voted.
func (p *Poll) AlreadyVoted(voterID string) bool {
	for _, voteEntry := range p.Votes {
		if voteEntry.UserID == voterID {
			return true
		}
	}

	return false
}

// IsExpired checks if the deadline has passed.
func (p *Poll) IsExpired() bool {
	return time.Now().After(p.Deadline)
}

// SubmitVote validates and stores a vote.
func (p *Poll) SubmitVote(v vote.Vote) error {

	// Validate vote

	if p.IsClosed {
		return errors.New("poll is closed")
	}

	if p.IsExpired() {
		return errors.New("poll has expired")
	}

	if p.AlreadyVoted(v.UserID) {
		return errors.New("you have already voted for this poll")
	}

	if !p.ValidateVoteCount(v.MovieIDs) {
		return errors.New("too many movies selected")
	}

	if p.HasDuplicateMovies(v) {
		return errors.New("duplicated votes for the same movie are not allowed")
	}

	for _, movieID := range v.MovieIDs {
		if !p.HasMovie(movieID) {
			return errors.New("this movie doesn't exist in this poll")
		}
	}

	p.AddVote(v)

	return nil
}
