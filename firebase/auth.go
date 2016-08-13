package firebase

import (
	"net/http"

	"github.com/dustin/gojson"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/pkg/errors"
)

const (
	googleIdentityURL          = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyAssertion?key=%s"
	googleSignUpURL            = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=%s"
	googlePasswordURL          = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
	googleEmailConfirmationURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getOobConfirmationCode?key=%s"
	googleSetAccountURL        = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/setAccountInfo?key=%s"
)

// Auth Firebase認証結果
type Auth struct {
	Token        string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	// User
	LocalID       string `json:"localId"`
	FederatedID   string `json:"federatedId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	DisplayName   string `json:"displayName"`
	ScreenName    string `json:"screenName"`
	PhotoURL      string `json:"photoUrl"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	ProviderID    string `json:"providerId"`
}

// AuthService Firebaseの認証を行うサービス
type AuthService struct {
	client *Client
}

// CreateUserWithEmailAndPassword emailとpasswordで新規ユーザ登録
func (s *AuthService) CreateUserWithEmailAndPassword(email, password string) (Auth, error) {
	return s.signIn(googleSignUpURL, map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
}

// SendPasswordResetEmail 指定のemailユーザにpasswordリセットを通知
func (s *AuthService) SendPasswordResetEmail(email string) error {
	return s.client.postNoResponse(googleEmailConfirmationURL, map[string]interface{}{
		"requestType": "PASSWORD_RESET",
		"email":       email,
	})
}

// SendEmailVerify 指定のemailユーザにメール確認を通知
func (s *AuthService) SendEmailVerify(idToken string) error {
	return s.client.postNoResponse(googleEmailConfirmationURL, map[string]interface{}{
		"requestType": "VERIFY_EMAIL",
		"idToken":     idToken,
	})
}

// SendNewEmailAccept 指定のemailユーザにメール確認を通知 TODO: Deprecated うまく動かない
func (s *AuthService) SendNewEmailAccept(idToken, oldEmail, nextEmail string) error {
	return s.client.postNoResponse(googleEmailConfirmationURL, map[string]interface{}{
		"requestType": "NEW_EMAIL_ACCEPT",
		"idToken":     idToken,
		"email":       oldEmail,
		"newEmail":    nextEmail,
	})
}

func (s *AuthService) SignInAnonymously() Auth {
	// TODO:
	/*
	   var content = $"{{\"returnSecureToken\":true}}";

	   return await this.SignInWithPostContentAsync(GoogleSignUpUrl, content).ConfigureAwait(false);
	*/
	return Auth{}
}

// SignInWithEmailAndPassword passwordとemailでsignInする
func (s *AuthService) SignInWithEmailAndPassword(email, password string) (Auth, error) {
	return s.signIn(googlePasswordURL, map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
}

// SignInWithOAuth OAuthProviderでログインする
func (s *AuthService) SignInWithOAuth(provider provider.Provider, postBody string) (Auth, error) {
	return s.signIn(googleIdentityURL, map[string]interface{}{
		"postBody":          postBody,
		"requestUri":        "http://localhost",
		"returnSecureToken": true,
	})
}

func (s *AuthService) LinkAccountsAsyncWithEmailAndPassword(auth Auth, email, password string) Auth {
	// TODO:
	/*
	   var content = $"{{\"idToken\":\"{auth.FirebaseToken}\",\"email\":\"{email}\",\"password\":\"{password}\",\"returnSecureToken\":true}}";

	   return await this.SignInWithPostContentAsync(GoogleSetAccountUrl, content).ConfigureAwait(false);

	*/
	return Auth{}
}

func (s *AuthService) LinkAccountsWithOAuth(auth Auth, provider provider.Provider, oauthAccessToken string) Auth {
	// TODO:
	/*
	   var providerId = this.GetProviderId(authType);
	   var content = $"{{\"idToken\":\"{auth.FirebaseToken}\",\"postBody\":\"access_token={oauthAccessToken}&providerId={providerId}\",\"requestUri\":\"http://localhost\",\"returnSecureToken\":true}}";

	   return await this.SignInWithPostContentAsync(GoogleIdentityUrl, content).ConfigureAwait(false);
	*/
	return Auth{}
}

func (s *AuthService) signIn(googleURL string, params map[string]interface{}) (Auth, error) {
	resp, err := s.client.post(googleURL, params)
	if err != nil {
		return Auth{}, errors.Wrapf(err, "request error params = %#v", params)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Auth{}, s.client.readBodyError(resp.StatusCode, resp.Body)
	}

	var auth Auth
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return Auth{}, errors.Wrap(err, "response json decode error")
	}

	return auth, nil
}
