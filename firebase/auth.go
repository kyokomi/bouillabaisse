package firebase

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dustin/gojson"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/pkg/errors"
)

const (
	googleIdentityURL      = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyAssertion?key=%s"
	googleSignUpURL        = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=%s"
	googlePasswordURL      = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
	googlePasswordResetURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getOobConfirmationCode?key=%s"
	googleSetAccountURL    = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/setAccountInfo?key=%s"
)

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
	PhotoUrl      string `json:"photoUrl"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	ProviderID    string `json:"providerId"`
}

// Config
// https://github.com/step-up-labs/firebase-authentication-dotnet/blob/master/src/Firebase.Auth/FirebaseConfig.cs

// AuthLink
// https://github.com/step-up-labs/firebase-authentication-dotnet/blob/master/src/Firebase.Auth/FirebaseAuthLink.cs

// User
// https://github.com/step-up-labs/firebase-authentication-dotnet/blob/master/src/Firebase.Auth/User.cs

type AuthService struct {
	client *Client
}

func (s *AuthService) CreateUserWithEmailAndPassword(email, password string) Auth {
	// TODO:
	/*
	   var content = $"{{\"email\":\"{email}\",\"password\":\"{password}\",\"returnSecureToken\":true}}";

	   return await this.SignInWithPostContentAsync(GoogleSignUpUrl, content).ConfigureAwait(false);
	*/
	return Auth{}
}

func (s *AuthService) SendPasswordResetEmailAsync(email string) error {
	// TODO:
	/*
	   var content = $"{{\"requestType\":\"PASSWORD_RESET\",\"email\":\"{email}\"}}";

	   var response = await this.client.PostAsync(new Uri(string.Format(GooglePasswordResetUrl, this.authConfig.ApiKey)), new StringContent(content, Encoding.UTF8, "application/json")).ConfigureAwait(false);

	   response.EnsureSuccessStatusCode();
	*/
	return nil
}

func (s *AuthService) SignInAnonymously() Auth {
	// TODO:
	/*
	   var content = $"{{\"returnSecureToken\":true}}";

	   return await this.SignInWithPostContentAsync(GoogleSignUpUrl, content).ConfigureAwait(false);
	*/
	return Auth{}
}

func (s *AuthService) SignInWithEmailAndPassword(email, password string) Auth {
	// TODO:
	/*
	   var content = $"{{\"email\":\"{email}\",\"password\":\"{password}\",\"returnSecureToken\":true}}";

	   return await this.SignInWithPostContentAsync(GooglePasswordUrl, content).ConfigureAwait(false);
	*/
	return Auth{}
}

func (s *AuthService) SignInWithOAuth(provider provider.Provider, postBody string) (Auth, error) {
	return s.signIn(googleIdentityURL, signInParams{
		PostBody:          postBody,
		RequestURI:        "http://localhost",
		ReturnSecureToken: true,
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

type signInParams struct {
	RequestURI        string `json:"requestUri"`
	PostBody          string `json:"postBody"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

func (s *AuthService) signIn(googleURL string, params signInParams) (Auth, error) {
	// Request Post
	body, err := json.Marshal(params)
	if err != nil {
		return Auth{}, errors.Wrapf(err, "signInParams Marshal error %#v", params)
	}
	url := fmt.Sprintf(googleURL, s.client.config.ApiKey)
	resp, err := s.client.httpClient.Post(
		url,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return Auth{}, errors.Wrapf(err, "%s request error params = %#v", url, params)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data = []byte{}
		}
		return Auth{}, errors.Errorf("response error statudCode = %d body = %s\n", resp.StatusCode, string(data))
	}

	var auth Auth
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return Auth{}, errors.Wrap(err, "response json decode error")
	}

	return auth, nil
}
