package firebase

import "net/http"

type Client struct {
	config     Config
	httpClient *http.Client

	Auth *AuthService
}

type Config struct {
	ApiKey string
}

func NewClient(config Config, transport http.RoundTripper) *Client {
	c := &Client{
		config:     config,
		httpClient: &http.Client{Transport: transport},
	}

	c.Auth = &AuthService{
		client: c,
	}

	return c
}
