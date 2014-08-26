package savey

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	defaultBaseURL   = "http://www.savey.co/"
	dashboardURL     = defaultBaseURL + "dashboard/"
	defaultUserAgent = "go-savey"
)

// Client represents savey client.
type Client struct {
	// http client
	client *http.Client

	// base url for all requests
	BaseURL *url.URL

	// user agent header
	UserAgent string

	// services
	Categories   *CategoryService
	Accounts     *AccountService
	Transactions *TransactionService
}

// CreateDefaultHTTPClient creates default HTTP with cookie jar.
func CreateDefaultHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &http.Client{Jar: jar}, nil
}

// NewClient creates new Savey client.
func NewClient(httpClient *http.Client) (*Client, error) {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	if httpClient == nil {
		httpClient, err = CreateDefaultHTTPClient()
		if err != nil {
			return nil, err
		}
	}

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
	c.Categories = &CategoryService{client: c}
	c.Accounts = &AccountService{client: c}
	c.Transactions = &TransactionService{client: c}

	return c, nil
}

// NewRequest creates a new HTTP request.
func (c *Client) NewRequest(method string, endpointURL string, body *url.Values) (*http.Request, error) {
	rel, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	buffer := new(bytes.Buffer)
	if body != nil {
		buffer = bytes.NewBufferString(body.Encode())
	}

	req, err := http.NewRequest(method, u.String(), buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

// Do performs HTTP request and returns HTTP response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Get creates a new GET HTTP request.
func (c *Client) Get(url string) (*http.Request, error) {
	return c.NewRequest("GET", url, nil)
}

// Post creates a new POST HTTP request.
func (c *Client) Post(url string, body *url.Values) (*http.Request, error) {
	return c.NewRequest("POST", url, body)
}

// Login logs in user with provided username and password.
func (c *Client) Login(username string, password string) error {
	payload := &url.Values{
		"identity":   {username},
		"credential": {password},
	}
	url := "user/login"
	req, err := c.Post(url, payload)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.Request.URL.String() != dashboardURL {
		return errors.New("Can not login with provided credentials.")
	}

	return nil
}

// Logout logs out currently logged user.
func (c *Client) Logout() error {
	url := "user/logout"
	req, err := c.Get(url)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.Request.URL.String() != defaultBaseURL {
		return errors.New("There was an error while logging out.")
	}

	return nil
}
