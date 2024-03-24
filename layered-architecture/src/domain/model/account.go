package model

// Account is a domain model that represents a bank account.
type Account struct {
	Id           int
	Balance      int
}

func NewAccount(balance int) *Account {
	return &Account{
		Balance:  balance,
	}
}
