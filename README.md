# Savey

Export data from your favourite money management tool http://www.savey.co/.

## Install

~~~
go get github.com/jbub/savey
~~~

## Example

~~~ go
package main

import (
	"encoding/json"
	"log"

	"github.com/jbub/savey"
)

const (
	userLogin    = "your@login.com"
	userPassword = "yourpassword"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := savey.NewClient(nil)
	checkErr(err)

	err = client.Login(userLogin, userPassword)
	checkErr(err)

	accounts, err := client.Accounts.GetAccounts()
	checkErr(err)

	categories, err := client.Categories.GetCategories()
	checkErr(err)

	transactions, err := client.Transactions.GetTransactions(accounts)
	checkErr(err)

	err = client.Logout()
	checkErr(err)
}
~~~
