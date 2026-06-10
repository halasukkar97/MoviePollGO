package api

import (
	"movie-vote/database"
	"movie-vote/movie"
	"movie-vote/poll"
	"movie-vote/user"
	"movie-vote/vote"
)

// FindPollByID searches the in-memory poll list for a poll with the given ID.
// It returns both the poll pointer and a boolean so callers can handle "not found".
func FindPollByID(pollID string) (*poll.Poll, bool) {
	for i := range polls {
		if polls[i].ID == pollID {
			return &polls[i], true
		}
	}

	return nil, false
}

// SavePoll stores a newly created poll in memory.
func SavePoll(poll poll.Poll) {
	polls = append(polls, poll)
}

// SaveMovie stores a newly created movie in memory.
func SaveMovie(movie movie.Movie) {
	movies = append(movies, movie)
}

// SaveUser stores a newly created user in PostgreSQL.
// Returning an error lets the HTTP handler send a clear failure response.
func SaveUser(user user.User) error {
	_, err := database.DB.Exec(
		"INSERT INTO users (id, name) VALUES ($1, $2)",
		user.ID,
		user.Name,
	)

	return err
}

// SaveVote stores a valid vote in memory after the poll accepts it.
func SaveVote(vote vote.Vote) {
	votes = append(votes, vote)
}

// GetAllPolls returns every poll currently stored in memory.
func GetAllPolls() []poll.Poll {
	return polls
}
