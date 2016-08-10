package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
)

type ServerContext struct {
	fireClient *firebase.Client
	config     Config
}

func (*ServerContext) loginHandler(c echo.Context) error {
	providerName := c.Param("provider")

	loginUrl, err := provider.GetBeginAuthURL(provider.New(providerName))
	if err != nil {
		return errors.Wrapf(err, "GetBeginAuthURLの呼び出し中にエラーが発生しました: %s", providerName)
	}

	return c.Redirect(http.StatusTemporaryRedirect, loginUrl)
}

func (s *ServerContext) callbackHandler(c echo.Context) error {
	providerName := c.Param("provider")

	p := provider.New(providerName)
	postBody, err := provider.BuildSignInPostBody(p, c.QueryParams())
	if err != nil {
		return errors.Wrapf(err, "%s BuildSignInPostBody error", p.Name())
	}

	auth, err := s.fireClient.Auth.SignInWithOAuth(p, postBody)
	if err != nil {
		return errors.Wrapf(err, "%s SignInWithOAuth error", p.Name())
	}

	pp.Println(auth) // TODO: debug

	now := time.Now()
	a := AuthStore{Auth: auth, CreatedAt: now, UpdateAt: now}
	stores.Add(a)

	return c.JSON(http.StatusOK, auth)
}

func ServeWithConfig(config Config) (string, error) {
	baseURL := fmt.Sprintf("http://localhost%s", config.Server.ListenAddr)

	// setup
	provider.InitOAuth(baseURL, config.Auth)

	s := ServerContext{
		fireClient: firebase.NewClient(
			firebase.Config{ApiKey: config.Server.FirebaseApiKey}, &http.Transport{},
		),
		config: config,
	}

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

	e.GET(provider.RESTAuthLoginPath(), s.loginHandler)
	e.GET(provider.RESTCallbackPath(), s.callbackHandler)

	// TODO: 適当にgoroutine
	go e.Run(standard.New(config.Server.ListenAddr))

	return baseURL, nil
}
