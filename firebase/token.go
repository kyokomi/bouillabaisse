package firebase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/stretchr/gomniauth/oauth2"
)

// https://developers.google.com/identity/toolkit/reference/securetoken/rest/v1/token
const (
	googleTokenURL = "https://securetoken.googleapis.com/v1/token?key=%s"
)

type Token struct {
	AccessToken  string `json:"access_token"` // The granted access token.
	ExpiresIn    string `json:"expires_in"`   // Expiration time of access_token in seconds.
	TokenType    string `json:"token_type"`   // The type of access_token. Included to conform with the OAuth 2.0 specification; always Bearer.
	RefreshToken string `json:"refresh_token"`
}

type TokenService struct {
	client *Client
}

func (s *TokenService) ExchangeRefreshToken(refreshToken string) (Token, error) {
	data := url.Values{}
	data.Set(oauth2.OAuth2KeyGrantType, oauth2.OAuth2KeyRefreshToken)
	data.Set(oauth2.OAuth2KeyRefreshToken, refreshToken)
	return s.exchangeToken(data)
}

func (s *TokenService) ExchangeAuthorizationCode(idToken string) (Token, error) {
	data := url.Values{}
	data.Set(oauth2.OAuth2KeyGrantType, oauth2.OAuth2GrantTypeAuthorizationCode)
	data.Set(oauth2.OAuth2KeyCode, idToken)
	return s.exchangeToken(data)
}

func (s *TokenService) exchangeToken(data url.Values) (Token, error) {
	// Request Post
	url := fmt.Sprintf(googleTokenURL, s.client.config.ApiKey)
	resp, err := s.client.httpClient.PostForm(url, data)
	if err != nil {
		return Token{}, errors.Wrapf(err, "%s request error params = %#v", url, data)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data = []byte{}
		}
		return Token{}, errors.Errorf("response error statudCode = %d body = %s\n", resp.StatusCode, string(data))
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return Token{}, errors.Wrap(err, "response json decode error")
	}

	return token, nil
}
