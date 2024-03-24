package main

import (
	"fmt"
	"log"

	"github.com/tkhrk1010/go-samples/go_bank_account/src/infrastructure/repository"
	"github.com/tkhrk1010/go-samples/go_bank_account/src/usecase"
)

func main() {
	repo := repository.NewInMemoryAccountRepository()
	usecase := usecase.NewAccountUsecase(repo)

	account, err := usecase.OpenAccount(1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Opened account with Id %d and balance %d\n", account.Id, account.Balance)

	account2, err := usecase.OpenAccount(2000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Opened account with Id %d and balance %d\n", account2.Id, account2.Balance)

	account, err = usecase.Deposit(account.Id, 500)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deposited 500 into account with Id %d; new balance is %d\n", account.Id, account.Balance)

	account, err = usecase.Withdraw(account.Id, 1500)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Withdrew 1500 from account with Id %d; new balance is %d\n", account.Id, account.Balance)

	err = usecase.CloseAccount(account.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("account with Id %d was Closed \n", account.Id)

	balance2, err := usecase.CheckBalance(account2.Id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance of account with Id %d is %d\n", account2.Id, balance2)

}
