package main

import (
	"fmt"
	"log"
	"movie-vote/api"
	"movie-vote/database"
	"net/http"

	"github.com/joho/godotenv"
)

// main is the first function Go runs when the application starts.
// It loads environment variables, connects to the database, registers every
// HTTP route, and then starts listening for requests on port 8080.
func main() {

	// godotenv.Load reads key/value pairs from .env so os.Getenv can use them.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	err = database.Connect()
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
	// /movies/search calls TMDB instead of the local movie database.
	http.HandleFunc("/movies/search", api.SearchMoviesHandler)
	// The trailing slash route catches paths like /polls/abc-123.
	http.HandleFunc("/polls/", api.PollByIDHandler)

	// ListenAndServe keeps the server running and waits for browser/API requests.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// MovieVoteHandler handles the root route and returns a simple health message.
func MovieVoteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Movie Vote API updated")
}
