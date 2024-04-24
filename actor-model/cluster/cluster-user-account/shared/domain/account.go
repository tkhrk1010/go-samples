package domain

import (
	"github.com/google/uuid"
)

type Account struct {
	ID    string
	Email string
}

func generateUUID() string {
	return uuid.New().String()
}

func NewAccount(email string) *Account {
	return &Account{
		ID: generateUUID(),
		Email: email,
	}
}
