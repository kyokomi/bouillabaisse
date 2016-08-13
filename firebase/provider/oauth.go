package provider

import (
	"fmt"

	"github.com/mrjones/oauth"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/oauth2"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"

	"github.com/kyokomi/bouillabaisse/twitter"
)

const (
	authLoginPath = "/auth/login"
	callbackPath  = "/auth/callback"
)

// BuildAuthLoginPath build login path
func BuildAuthLoginPath() string {
	return fmt.Sprintf("%s/:provider", authLoginPath)
}

// BuildCallbackPath build callback path
func BuildCallbackPath() string {
	return fmt.Sprintf("%s/:provider", callbackPath)
}

var twitterOAuthClient *twitter.OAuthClient

// InitOAuth OAuthプロバイダー各社の初期化を行います
func InitOAuth(baseURL string, config Config) {
	gomniauth.SetSecurityKey(config.AuthSecretKey)
	gomniauth.WithProviders(
		github.New(
			config.GitHubClientID,
			config.GitHubSecretKey,
			buildCallbackURL(GitHubProvider, baseURL),
		),
		google.New(
			config.GoogleClientID,
			config.GoogleSecretKey,
			buildCallbackURL(GoogleProvider, baseURL),
		),
		facebook.New(
			config.FacebookID,
			config.FacebookSecretKey,
			buildCallbackURL(FacebookProvider, baseURL),
		),
	)
	twitterOAuthClient = twitter.NewTwitterOAuth(
		config.TwitterConsumerID,
		config.TwitterConsumerSecretKey,
		buildCallbackURL(TwitterProvider, baseURL),
	)
}

func buildCallbackURL(p Provider, domain string) string {
	return fmt.Sprintf("%s%s/%s", domain, callbackPath, p.Name())
}

// BuildSignInURL 指定providerのsignInURLをbuildします
func BuildSignInURL(p Provider, domain string) string {
	return fmt.Sprintf("%s%s/%s", domain, authLoginPath, p.Name())
}

// GetBeginAuthURL 指定したproviderの認証URLを取得します
func GetBeginAuthURL(p Provider) (string, error) {
	if p == TwitterProvider {
		return twitterOAuthClient.GetRequestTokenAndURL()
	}

	provider, err := gomniauth.Provider(p.Name())
	if err != nil {
		return "", err
	}
	return provider.GetBeginAuthURL(nil, nil)
}

// BuildSignInPostBody 指定したproviderの認証に必要なpostBodyを生成します
func BuildSignInPostBody(p Provider, params map[string][]string) (string, error) {
	return providerPostBodyFuncMap[p](p, params)
}

type providerPostBodyFunc func(p Provider, params map[string][]string) (string, error)

var providerPostBodyFuncMap = map[Provider]providerPostBodyFunc{
	GitHubProvider:   oauth2ProviderPostBody,
	TwitterProvider:  twitterProviderPostBody,
	FacebookProvider: oauth2ProviderPostBody,
	GoogleProvider:   oauth2ProviderPostBody,
}

func oauth2ProviderPostBody(p Provider, params map[string][]string) (string, error) {
	creds, err := completeAuth(p, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s=%s&providerId=%s",
		oauth2.OAuth2KeyAccessToken, creds.Get(oauth2.OAuth2KeyAccessToken).String(), p.id()), nil
}

func completeAuth(p Provider, params map[string][]string) (*common.Credentials, error) {
	var rawQuery string
	for k, v := range params {
		if len(v) <= 0 {
			continue
		}
		rawQuery += fmt.Sprintf("&%s=%s", k, v[0])
	}
	rawQuery = rawQuery[1:] // delete start &

	provider, err := gomniauth.Provider(p.Name())
	if err != nil {
		return nil, err
	}

	return provider.CompleteAuth(objx.MustFromURLQuery(rawQuery))
}

func twitterProviderPostBody(p Provider, params map[string][]string) (string, error) {
	verificationCode := params[oauth.VERIFIER_PARAM][0]
	tokenKey := params[oauth.TOKEN_PARAM][0]

	accessToken, err := twitterOAuthClient.GetAccessToken(tokenKey, verificationCode)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s=%s&%s=%s&providerId=%s",
		oauth2.OAuth2KeyAccessToken, accessToken.Token,
		oauth.TOKEN_SECRET_PARAM, accessToken.Secret, p.id()), nil
}
