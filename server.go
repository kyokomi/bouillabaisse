package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"gopkg.in/go-pp/pp.v2"
	"gopkg.in/yaml.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/labstack/gommon/log"
)

type ServerContext struct {
	fireClient *firebase.Client
	config     ServerConfig
}

type ServerConfig struct {
	ListenAddr        string
	FirebaseApiKey    string
	AuthStoreDirPath  string
	AuthStoreFileName string

	AuthConfig provider.Config
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
	if err := a.Save(s.config.AuthStoreDirPath, s.config.AuthStoreFileName); err != nil {
		log.Errorf("%+v", errors.Wrapf(err, "%s StoreSave error", p.Name()))
		// 後続処理は行う
	}

	return c.JSON(http.StatusOK, auth)
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

	s := ServerContext{
		fireClient: firebase.NewClient(
			firebase.Config{ApiKey: config.FirebaseApiKey}, &http.Transport{},
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
	go e.Run(standard.New(config.ListenAddr))

	return baseURL, nil
}
