# Savey

[![Build Status](https://travis-ci.org/jbub/savey.svg)](https://travis-ci.org/jbub/savey)

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
    // create new savey client with default HTTP client
	client, err := savey.NewClient(nil)
	checkErr(err)
    
    // login user using username and password
	err = client.Login(userLogin, userPassword)
	checkErr(err)
    
    // list payment accounts of current user
	accounts, err := client.Accounts.GetAccounts()
	checkErr(err)
    
    // list categories of current user
	categories, err := client.Categories.GetCategories()
	checkErr(err)
    
    // list all transactions for given accounts of current user
	transactions, err := client.Transactions.GetTransactions(accounts)
	checkErr(err)
    
    // logout user
	err = client.Logout()
	checkErr(err)
}
~~~
