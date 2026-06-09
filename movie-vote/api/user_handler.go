package api

import (
	"encoding/json"
	"movie-vote/user"
	"net/http"
)

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name string `json:"name"`
}

// CreateUserResponse is returned after a user is created.
type CreateUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateUserHandler handles POST /users.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	createdUser := user.CreateNewUser(user.CreateUserInput{
		Name: req.Name,
	})

	users = append(users, createdUser)

	response := CreateUserResponse{
		ID:   createdUser.ID,
		Name: createdUser.Name,
	}

	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)

		return
	}
}

// UsersHandler routes requests to the correct user handler.
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CreateUserHandler(w, r)
		return
	}

	if r.Method == http.MethodGet {
		ListUsersHandler(w, r)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

// ListUsersHandler handles GET /users.
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
