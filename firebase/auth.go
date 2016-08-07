package firebase

import (
	"bytes"
	"fmt"

	"net/http"

	"io/ioutil"

	"github.com/bitly/go-simplejson"
	"github.com/dustin/gojson"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/pkg/errors"
	"gopkg.in/go-pp/pp.v2"
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
	ExpiresIn    int    `json:"expiresIn"`
	// User
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
	return Auth{}, s.signInWithPostContent(googleIdentityURL, signInParams{
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
	PostBody          string `json:"postBody"`
	RequestURI        string `json:"requestUri"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
	SessionID         string `json:"sessionId"`
}

func (s *AuthService) signInWithPostContent(googleURL string, params signInParams) error {
	// Request Post
	body, err := json.Marshal(params)
	if err != nil {
		return errors.Wrapf(err, "signInParams Marshal error %#v", params)
	}
	url := fmt.Sprintf(googleURL, s.client.config.ApiKey)
	pp.Println(url, string(body)) // TODO: debug
	resp, err := s.client.httpClient.Post(
		url,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.Wrapf(err, "%s request error params = %#v", googleURL, params)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data = []byte{}
		}
		return errors.Errorf("response error statudCode = %d body = %s\n", resp.StatusCode, string(data))
	}

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "%s response Marshal error")
	}
	pp.Println(json)

	return nil
}
