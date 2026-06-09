package api

import (
	"movie-vote/movie"
	"movie-vote/poll"
	"movie-vote/user"
	"movie-vote/vote"
)

// FindPollByID returns the matching poll and whether it was found.
func FindPollByID(pollID string) (*poll.Poll, bool) {
	for i := range polls {
		if polls[i].ID == pollID {
			return &polls[i], true
		}
	}

	return nil, false
}

func SavePoll(poll poll.Poll) {
	polls = append(polls, poll)
}

func SaveMovie(movie movie.Movie) {
	movies = append(movies, movie)
}

func SaveUser(user user.User) {
	users = append(users, user)
}

func SaveVote(vote vote.Vote) {
	votes = append(votes, vote)
}

func GetAllPolls() []poll.Poll {
	return polls
}
