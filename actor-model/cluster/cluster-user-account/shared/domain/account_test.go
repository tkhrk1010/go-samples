package domain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
)

func TestNewAccount(t *testing.T) {
	id := "user123"
	email := "user@example.com"

	account := domain.NewAccount(id, email)

	if account.ID != id {
		t.Errorf("Expected account ID to be %s, but got %s", id, account.ID)
	}

	if account.Email != email {
		t.Errorf("Expected account email to be %s, but got %s", email, account.Email)
	}
}