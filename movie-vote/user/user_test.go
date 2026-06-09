package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateNewUser(t *testing.T) {
	input := CreateUserInput{Name: "Hela"}

	u := CreateNewUser(input)

	if _, err := uuid.Parse(u.ID); err != nil {
		t.Fatalf("expected valid UUID, got %q: %v", u.ID, err)
	}

	if u.Name != input.Name {
		t.Errorf("expected name %q, got %q", input.Name, u.Name)
	}
}
