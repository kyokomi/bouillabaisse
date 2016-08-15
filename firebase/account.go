package firebase

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	googleGetAccountURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo?key=%s"
)

// AccountInfo firebaseで管理しているAccountの情報
type AccountInfo struct {
	Kind  string `json:"kind"`
	Users []User `json:"users"`
}

// User Accountに紐づくUser
type User struct {
	LocalID           string     `json:"localId"`
	DisplayName       string     `json:"displayName"`
	PhotoURL          string     `json:"photoUrl"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"passwordHash"`
	EmailVerified     bool       `json:"emailVerified"`
	ValidSince        string     `json:"validSince"`
	ProviderUserInfo  []UserInfo `json:"providerUserInfo"`
	LastLoginAt       string     `json:"lastLoginAt"`
	PasswordUpdatedAt float64    `json:"passwordUpdatedAt"`
	CreatedAt         string     `json:"createdAt"`
}

// UserInfo Userに紐づく書くProviderの情報
type UserInfo struct {
	RawID       string `json:"rawId"`
	ProviderID  string `json:"providerId"`
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoUrl"`
	Email       string `json:"email"`
	FederatedID string `json:"federatedId"`
}

// AccountService is account manager service
type AccountService struct {
	client *Client
}

// GetAccountInfo return AccountInfo by idToken
func (s *AccountService) GetAccountInfo(idToken string) (AccountInfo, error) {
	params := map[string]interface{}{
		"idToken": idToken,
	}
	resp, err := s.client.post(googleGetAccountURL, params)
	if err != nil {
		return AccountInfo{}, errors.Wrapf(err, "request error params = %#v", params)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AccountInfo{}, s.client.readBodyError(resp.StatusCode, resp.Body)
	}

	var accountInfo AccountInfo
	if err := json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		return AccountInfo{}, errors.Wrap(err, "response json decode error")
	}

	return accountInfo, nil
}
