package user

import "github.com/google/uuid"

type User struct {
	ID   string
	Name string
}

type CreateUserInput struct {
	Name string
}

func CreateNewUser(input CreateUserInput) User {
	return User{
		ID:   uuid.New().String(),
		Name: input.Name,
	}
}
