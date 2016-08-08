package firebase

import "net/http"

type Client struct {
	config     Config
	httpClient *http.Client

	Auth  *AuthService
	Token *TokenService
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
	c.Token = &TokenService{
		client: c,
	}

	return c
}
