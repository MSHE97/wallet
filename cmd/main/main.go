package main

import (
	"fmt"
	"github.com/MSHE97/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}
		return
	}

	_, err = svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = svc.RegisterAccount("+992000000003")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.ExportToFile("export.txt")
	if err != nil {
		return
	}
	err = svc.ImportFromFile("import.txt")
	if err != nil {
		return
	}
	importAccount, err := svc.FindAccountById(6)
	if err != nil {
		log.Print(err)
	}
	log.Print(*importAccount)

	fmt.Println(account.ID)			// 1
	fmt.Println(account.Balance)	// 10
}
