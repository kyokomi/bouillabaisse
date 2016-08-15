package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/go-pp/pp.v2"
)

// Client is firebase client
type Client struct {
	config     Config
	httpClient *http.Client

	Auth    *AuthService
	Token   *TokenService
	Account *AccountService
}

func (c *Client) postNoResponse(googleURL string, params map[string]interface{}) error {
	resp, err := c.post(googleURL, params)
	if err != nil {
		return errors.Wrapf(err, "request error params = %#v", params)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.readBodyError(resp.StatusCode, resp.Body)
	}
	return nil
}

func (c *Client) readBodyError(statusCode int, body io.ReadCloser) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		data = []byte{}
	}
	return errors.Errorf("response error statudCode = %d body = %s\n", statusCode, string(data))
}

func (c *Client) post(googleURL string, params map[string]interface{}) (*http.Response, error) {
	pp.Println(params) // debug log
	// Request Post
	body, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrapf(err, "params Marshal error %#v", params)
	}
	url := fmt.Sprintf(googleURL, c.config.APIKey)
	return c.httpClient.Post(
		url,
		"application/json",
		bytes.NewReader(body),
	)
}

// Config firebase client configuration
type Config struct {
	APIKey string
}

// NewClient create firebase client
func NewClient(config Config, transport http.RoundTripper) *Client {
	c := &Client{
		config:     config,
		httpClient: &http.Client{Transport: transport},
	}

	c.Auth = &AuthService{
		client: c,
	}
	c.Token = &TokenService{
		client: c,
	}
	c.Account = &AccountService{
		client: c,
	}

	return c
}
