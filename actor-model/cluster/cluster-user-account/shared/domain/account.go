package domain

type Account struct {
	ID    string
	Email string
}


func NewAccount(id string, email string) *Account {
	return &Account{
		ID: id,
		Email: email,
	}
}
