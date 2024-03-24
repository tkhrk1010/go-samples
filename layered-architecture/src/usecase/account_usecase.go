package usecase

import (
	"errors"

	"github.com/tkhrk1010/go-samples/go_bank_account/src/domain/model"
	"github.com/tkhrk1010/go-samples/go_bank_account/src/domain/repository"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type AccountUsecase interface {
	OpenAccount(initialDeposit int) (*model.Account, error)
	CheckBalance(accountId int) (int, error)
	Deposit(accountId int, amount int) (*model.Account, error)
	Withdraw(accountId int, amount int) (*model.Account, error)
	CloseAccount(accountId int) error
}

type accountUsecase struct {
	repo repository.AccountRepository
}

func NewAccountUsecase(repo repository.AccountRepository) AccountUsecase {
	return &accountUsecase{
		repo: repo,
	}
}

func (a *accountUsecase) OpenAccount(initialDeposit int) (*model.Account, error) {
	if initialDeposit < 0 {
		return nil, errors.New("initial deposit must be positive")
	}
	newAccount := model.NewAccount(initialDeposit)
	err := a.repo.Save(newAccount)
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}

func (a *accountUsecase) CheckBalance(accountId int) (int, error) {
	account, err := a.repo.FindById(accountId)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (a *accountUsecase) Deposit(accountId int, amount int) (*model.Account, error) {
	account, err := a.repo.FindById(accountId)
	if err != nil {
		return nil, err
	}

	if amount <= 0 {
		return nil, errors.New("invalid deposit amount")
	}

	account.Balance += amount
	err = a.repo.Update(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountUsecase) Withdraw(accountId int, amount int) (*model.Account, error) {
	// Retrieve account from repository
	account, err := a.repo.FindById(accountId)
	if err != nil {
		return nil, err
	}

	// Check if there are sufficient funds
	if account.Balance < amount {
		return nil, ErrInsufficientFunds
	}

	// Deduct amount from balance
	account.Balance -= amount

	// Update account in repository
	err = a.repo.Update(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountUsecase) CloseAccount(accountId int) error {
	// Accountを検索します
	account, err := a.repo.FindById(accountId)
	if err != nil {
		return err
	}

	// Accountを削除します
	err = a.repo.Delete(account)
	if err != nil {
		return err
	}

	return nil
}
