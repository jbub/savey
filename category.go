package savey

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// CategoryService provides methods for working with categories.
type CategoryService struct {
	client *Client
}

// Category represents transaction category.
type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Kind  string `json:"kind"`
}

// GetCategories lists categories for current user.
func (s *CategoryService) GetCategories() ([]Category, error) {
	url := "manage"
	req, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	return ParseCategories(resp)
}

// ParseCategories parses categories from HTTP response.
func ParseCategories(resp *http.Response) ([]Category, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	categories := []Category{}
	var catErr interface{}
	doc.Find(".setup-categories .section.group").Not(".setup-heading").EachWithBreak(func(i int, s *goquery.Selection) {
		cat, err := ParseCategory(s)
		if err != nil {
			catErr = err
			return false
		}
		categories = append(categories, *cat)
	})
	if catErr != nil {
		return nil, catErr.(error)
	}
	return categories, nil
}

// ParseCategory parses category from selection.
func ParseCategory(s *goquery.Selection) (*Category, error) {
	title := s.Find(".col.span_6_of_12").Text()
	kind := s.Find(".col.span_2_of_12").Text()
	attr, exists := s.Find(".col.span_3_of_12 a").Eq(0).Attr("onclick")
	if !exists {
		return nil, errors.New("Category id attr not found.")
	}
	id, err := ParseID(attr)
	if err != nil {
		return nil, err
	}
	c := &Category{
		ID:    id,
		Title: CleanString(title),
		Kind:  CleanString(kind),
	}
	return c, nil
}
