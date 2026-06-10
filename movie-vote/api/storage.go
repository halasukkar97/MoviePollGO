package api

import (
	"movie-vote/movie"
	"movie-vote/vote"
)

// These slices act as simple in-memory storage for parts of the app that have
// not been moved to PostgreSQL yet. Data here is lost when the server restarts.
var movies []movie.Movie
var votes []vote.Vote
