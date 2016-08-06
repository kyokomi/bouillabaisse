package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"gopkg.in/go-pp/pp.v2"
	"gopkg.in/yaml.v2"

	"github.com/kyokomi/bouillabaisse/twitter"
)

var twitterOAuthClient *twitter.OAuthClient

type ServerConfig struct {
	ListenAddr               string
	AuthSecretKey            string // gomniauth setup secretKey
	GitHubClientID           string
	GitHubSecretKey          string
	GoogleClientID           string
	GoogleSecretKey          string
	FacebookID               string
	FacebookSecretKey        string
	TwitterConsumerID        string
	TwitterConsumerSecretKey string
}

func loginHandler(c echo.Context) error {
	providerName := c.Param("provider")

	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		if providerName == TwitterProvider.Name() {
			return loginTwitterHandler(c)
		}
		return errors.Wrapf(err, "認証プロバイダーの取得に失敗しました %s", provider)
	}

	loginUrl, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		return errors.Wrapf(err, "GetBeginAuthURLの呼び出し中にエラーが発生しました: %s", provider)
	}

	return c.Redirect(http.StatusTemporaryRedirect, loginUrl)
}

func loginTwitterHandler(c echo.Context) error {
	requestURL, err := twitterOAuthClient.GetRequestTokenAndURL()
	if err != nil {
		return errors.Wrap(err, "GetRequestTokenAndURL error")
	}

	return c.Redirect(http.StatusFound, requestURL)
}

func callbackHandler(c echo.Context) error {
	providerName := c.Param("provider")

	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		if providerName == TwitterProvider.Name() {
			return authTwitterCallbackHandler(c)
		}
		return errors.Wrapf(err, "認証プロバイダーの取得に失敗しました %s", provider)
	}

	rawQuery := c.Request().URL().QueryString()

	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(rawQuery))
	if err != nil {
		return errors.Wrapf(err, "認証を完了できませんでした %s", provider)
	}

	return c.JSON(http.StatusOK, creds)
}

func authTwitterCallbackHandler(c echo.Context) error {
	verificationCode := c.QueryParam("oauth_verifier")
	tokenKey := c.QueryParam("oauth_token")

	accessToken, err := twitterOAuthClient.GetAccessToken(tokenKey, verificationCode)
	if err != nil {
		return errors.Wrap(err, "GetAccessToken error")
	}

	return c.JSON(http.StatusOK, accessToken)
}

func Serve(env, configPath string) (string, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	var cnf map[string]ServerConfig
	if err := yaml.Unmarshal(buf, &cnf); err != nil {
		return "", err
	}

	pp.Println("config => ", cnf)

	return ServeConfig(cnf[env])
}

func ServeConfig(config ServerConfig) (string, error) {
	baseURL := fmt.Sprintf("http://localhost%s", config.ListenAddr)

	// setup
	gomniauth.SetSecurityKey(config.AuthSecretKey)
	gomniauth.WithProviders(
		github.New(
			config.GitHubClientID,
			config.GitHubSecretKey,
			GitHubProvider.CallbackURL(baseURL),
		),
		google.New(
			config.GoogleClientID,
			config.GoogleSecretKey,
			GoogleProvider.CallbackURL(baseURL),
		),
		facebook.New(
			config.FacebookID,
			config.FacebookSecretKey,
			FacebookProvider.CallbackURL(baseURL),
		),
	)
	twitterOAuthClient = twitter.NewTwitterOAuth(
		config.TwitterConsumerID,
		config.TwitterConsumerSecretKey,
		TwitterProvider.CallbackURL(baseURL),
	)

	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		return echo.ErrNotFound
	})

	e.GET(authLoginPath+"/:provider", loginHandler)
	e.GET(callbackPath+"/:provider", callbackHandler)

	// TODO: 適当にgoroutine
	go e.Run(standard.New(config.ListenAddr))

	return baseURL, nil
}
