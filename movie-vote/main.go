package main

import (
	"fmt"
	"log"
	"movie-vote/api"
	"movie-vote/database"
	"net/http"
)

// main is the first function Go runs when the application starts.
// It connects to the database, registers every HTTP route, and then starts
// listening for requests on port 8080.
func main() {

	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected!")

	// http.HandleFunc connects a URL path to the function that should handle it.
	http.HandleFunc("/", MovieVoteHandler)
	http.HandleFunc("/polls", api.PollsHandler)
	http.HandleFunc("/users", api.UsersHandler)
	http.HandleFunc("/movies", api.MoviesHandler)
	http.HandleFunc("/votes", api.CreateVoteHandler)
	http.HandleFunc("/results", api.ResultsHandler)

	// ListenAndServe keeps the server running and waits for browser/API requests.
	http.ListenAndServe(":8080", nil)
}

// MovieVoteHandler handles the root route and returns a simple health message.
func MovieVoteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Movie Vote API updated")
}
