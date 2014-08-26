package savey

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// TransactionService provides methods for working with payment transactions.
type TransactionService struct {
	client *Client
}

// Transaction represents payment transaction.
type Transaction struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	CategoryID int64     `json:"category_id"`
	AccountID  int64     `json:"account_id"`
	Value      float64   `json:"value"`
	Currency   string    `json:"currency"`
	Expense    bool      `json:"expense"`
}

// GetTransactions lists transactions for given accounts.
func (s *TransactionService) GetTransactions(accounts []Account) ([]Transaction, error) {
	transactions := []Transaction{}
	for _, account := range accounts {
		accTransactions, err := s.GetAccountTransactions(account)
		if err != nil {
			return nil, err
		}
		for _, transaction := range accTransactions {
			transactions = append(transactions, transaction)
		}
	}
	return transactions, nil
}

// GetAccountTransactions lists transactions for given account.
func (s *TransactionService) GetAccountTransactions(acc Account) ([]Transaction, error) {
	url := fmt.Sprintf("/set-account/%v", acc.ID)
	req, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Request.URL.String() != dashboardURL {
		return nil, errors.New("Failed to switch account.")
	}
	return ParseTransactions(acc, resp)
}

// ParseTransactions parses account transactions from HTTP response.
func ParseTransactions(acc Account, resp *http.Response) ([]Transaction, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	transactions := []Transaction{}
	doc.Find(".section.group.list").Each(func(i int, s *goquery.Selection) {
		trans, err := ParseTransaction(acc, s)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, *trans)
	})
	return transactions, nil
}

// ParseTransaction parses account transaction from selection.
func ParseTransaction(acc Account, s *goquery.Selection) (*Transaction, error) {
	s1 := s.Find(".col.span_2_of_3.list-description")
	txt := s1.Find(".list-pad-left").Text()
	title := strings.Split(txt, "\n")[1]
	spans := s1.Find("span")
	date, err := ParseDate(CleanString(spans.Eq(0).Text()))
	if err != nil {
		return nil, err
	}
	attr1, exists1 := spans.Eq(1).Find("a").Attr("href")
	if !exists1 {
		return nil, errors.New("Transaction attr href not found.")
	}
	categoryID, err := ParseID(attr1)
	if err != nil {
		return nil, err
	}
	attr2, exists2 := spans.Eq(2).Find("a.edit-list").Attr("onclick")
	if !exists2 {
		return nil, errors.New("Transaction attr onclick not found.")
	}
	transactionID, err := ParseID(attr2)
	if err != nil {
		return nil, err
	}
	s2 := s.Find(".col.span_1_of_3.list-amount .list-pad-right span").Eq(0).Text()
	value, err := strconv.ParseFloat(s2, 64)
	if err != nil {
		return nil, err
	}
	t := &Transaction{
		ID:         transactionID,
		Title:      CleanString(title),
		Date:       date,
		CategoryID: categoryID,
		AccountID:  acc.ID,
		Value:      value,
		Currency:   acc.Currency,
		Expense:    value < 0,
	}
	return t, nil
}