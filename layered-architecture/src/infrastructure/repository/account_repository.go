package repository

import (
	"errors"

	"github.com/tkhrk1010/go-samples/go_bank_account/src/domain/model"
)

// implementation of the AccountRepository interface
type InMemoryAccountRepository struct {
	accounts map[int]*model.Account
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[int]*model.Account),
	}
}

func (r *InMemoryAccountRepository) Save(a *model.Account) error {
	if a.Id == 0 {
		a.Id = len(r.accounts) + 1
	}
	r.accounts[a.Id] = a
	return nil
}

func (r *InMemoryAccountRepository) FindById(id int) (*model.Account, error) {
	a, ok := r.accounts[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return a, nil
}

func (r *InMemoryAccountRepository) Update(a *model.Account) error {
	for i, account := range r.accounts {
		if account.Id == a.Id {
			r.accounts[i] = a
			return nil
		}
	}
	return ErrAccountNotFound
}

func (r *InMemoryAccountRepository) Delete(a *model.Account) error {
	for _, account := range r.accounts {
		if account.Id == a.Id {
			delete(r.accounts, a.Id)
			return nil
		}
	}
	return ErrAccountNotFound
}

var ErrAccountNotFound = errors.New("account not found")
