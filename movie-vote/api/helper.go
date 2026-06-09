package api

import (
	"movie-vote/poll"
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
