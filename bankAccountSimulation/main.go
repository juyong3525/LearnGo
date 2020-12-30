package main

import (
	"fmt"

	"github.com/yong/bankAndDictionaryProject/accounts"
)

func main() {
	account := accounts.NewAccount("yong")
	account.Deposit(10)
	fmt.Println(account)
}
