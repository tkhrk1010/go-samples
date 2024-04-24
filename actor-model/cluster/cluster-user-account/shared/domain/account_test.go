package domain_test

import (
	"testing"
	"github.com/google/uuid"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
)

func TestNewAccount(t *testing.T) {
	email := "user@example.com"

	account := domain.NewAccount(email)

	// uuidの生成をdomain層に置くのがどうなのかは微妙なライン。ここではこだわらずいいことにする。
	if _, err := uuid.Parse(account.ID); err != nil {
		t.Errorf("Expected account ID to be a valid UUID, but got %s", account.ID)
	}

	if account.Email != email {
		t.Errorf("Expected account email to be %s, but got %s", email, account.Email)
	}
}