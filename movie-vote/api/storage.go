package api

import (
	"movie-vote/movie"
	"movie-vote/poll"
	"movie-vote/user"
	"movie-vote/vote"
)

// These slices act as simple in-memory storage for parts of the app.
// Data saved here is lost when the server restarts; users are now stored in PostgreSQL.
var users []user.User
var polls []poll.Poll
var movies []movie.Movie
var votes []vote.Vote
