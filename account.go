package savey

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// AccountService provides methods for working with payment accounts.
type AccountService struct {
	client *Client
}

// Account represents payment account.
type Account struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Currency string `json:"currency"`
}

// GetAccounts lists accounts for current user.
func (s *AccountService) GetAccounts() ([]Account, error) {
	url := "manage"
	req, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	return ParseAccounts(resp)
}

// ParseAccounts parses accounts from HTTP response.
func ParseAccounts(resp *http.Response) ([]Account, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	var accErr interface{}
	doc.Find(".setup-accounts .section.group").Not(".setup-heading").EachWithBreak(func(i int, s *goquery.Selection) bool {
		acc, err := ParseAccount(s)
		if err != nil {
			accErr = err
			return false
		}
		accounts = append(accounts, *acc)
		return true
	})
	if accErr != nil {
		return nil, accErr.(error)
	}
	return accounts, nil
}

// ParseAccount parses account from selection.
func ParseAccount(s *goquery.Selection) (*Account, error) {
	title := s.Find(".col.span_6_of_12").Text()
	currency := s.Find(".col.span_2_of_12").Text()
	attr, exists := s.Find(".col.span_3_of_12 a").Eq(0).Attr("onclick")
	if !exists {
		return nil, errors.New("Account id attr not found.")
	}
	id, err := ParseID(attr)
	if err != nil {
		return nil, err
	}
	acc := &Account{
		ID:       id,
		Title:    CleanString(title),
		Currency: CleanString(currency),
	}
	return acc, nil
}
