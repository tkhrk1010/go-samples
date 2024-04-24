package domain_test

import (
	"github.com/google/uuid"
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
)

func TestNewAccount(t *testing.T) {
	id := "123e4567-e89b-12d3-a456-426614174000"
	email := "user@example.com"

	account := domain.NewAccount(id, email)

	// uuidの生成をdomain層に置くのがどうなのかは微妙なライン。ここではこだわらずいいことにする。
	if _, err := uuid.Parse(account.ID); err != nil {
		t.Errorf("Expected account ID to be a valid UUID, but got %s", account.ID)
	}

	if account.Email != email {
		t.Errorf("Expected account email to be %s, but got %s", email, account.Email)
	}
}
