package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"gopkg.in/go-pp/pp.v2"
	"gopkg.in/yaml.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
)

var fireClient *firebase.Client

type ServerConfig struct {
	ListenAddr     string
	FirebaseApiKey string

	AuthConfig provider.Config
}

func loginHandler(c echo.Context) error {
	providerName := c.Param("provider")

	loginUrl, err := provider.GetBeginAuthURL(provider.New(providerName))
	if err != nil {
		return errors.Wrapf(err, "GetBeginAuthURLの呼び出し中にエラーが発生しました: %s", providerName)
	}

	return c.Redirect(http.StatusTemporaryRedirect, loginUrl)
}

func callbackHandler(c echo.Context) error {
	providerName := c.Param("provider")

	p := provider.New(providerName)
	postBody, err := provider.BuildSignInPostBody(p, c.QueryParams())
	if err != nil {
		return errors.Wrapf(err, "%s BuildSignInPostBody error", p.Name())
	}
	if _, err := fireClient.Auth.SignInWithOAuth(p, postBody); err != nil {
		return errors.Wrapf(err, "%s SignInWithOAuth error", p.Name())
	}

	return c.JSON(http.StatusOK, `{"message": "TODO: OK"`)
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
	provider.InitOAuth(baseURL, config.AuthConfig)

	fireClient = firebase.NewClient(firebase.Config{ApiKey: config.FirebaseApiKey}, &http.Transport{})

	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.SetHTTPErrorHandler(func(err error, ctx echo.Context) {
		fmt.Printf("%+v", err)
		e.DefaultHTTPErrorHandler(err, ctx)
	})

	// handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		return echo.ErrNotFound
	})

	e.GET(provider.RESTAuthLoginPath(), loginHandler)
	e.GET(provider.RESTCallbackPath(), callbackHandler)

	// TODO: 適当にgoroutine
	go e.Run(standard.New(config.ListenAddr))

	return baseURL, nil
}
