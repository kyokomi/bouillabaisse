package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/rcrowley/goagain"

	"github.com/kyokomi/bouillabaisse/config"
	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
)

// AuthCallbackListenerFunc 認証後に呼ばれるcallbackのWebAPIの処理のListener
type AuthCallbackListenerFunc func(ctx echo.Context) error

// ProviderServeWithConfig serve provider auth server
func ProviderServeWithConfig(p provider.Provider, c config.Config, fn AuthCallbackListenerFunc) error {
	baseURL := fmt.Sprintf("http://localhost%s", c.Server.ListenAddr)

	// setup
	provider.InitOAuth(baseURL, c.Auth)

	e := createEcho(serverContext{
		fireClient: firebase.NewClient(
			firebase.Config{APIKey: c.Server.FirebaseAPIKey}, &http.Transport{},
		),
		config: c,
	}, fn)
	l, err := createGoAgainListener(e, c)
	if err != nil {
		return err
	}

	signInURL := provider.BuildSignInURL(p, baseURL)
	fmt.Fprintln(os.Stdout, fmt.Sprintf("server stopped [ctrl + c] \n\n signInURL => %s", signInURL))

	return goagaginWait(l)
}

func createEcho(s serverContext, fn AuthCallbackListenerFunc) *echo.Echo {
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

	e.GET(provider.BuildAuthLoginPath(), s.loginHandler)
	e.GET(provider.BuildCallbackPath(), func(ctx echo.Context) error {
		return s.callbackHandler(ctx, fn)
	})

	return e
}

func createGoAgainListener(e *echo.Echo, c config.Config) (net.Listener, error) {
	l, err := goagain.Listener()
	if nil != err {
		// Listen on a TCP or a UNIX domain socket (TCP here).
		l, err = net.Listen("tcp", fmt.Sprintf("localhost%s", c.Server.ListenAddr))
		if nil != err {
			return nil, err
		}
		log.Println("listening on", l.Addr())

		// Accept connections in a new goroutine.
		go serve(e, l)
	} else {
		// Resume accepting connections in a new goroutine.
		log.Println("resuming listening on", l.Addr())
		go serve(e, l)

		// Kill the parent, now that the child has started successfully.
		if err := goagain.Kill(); err != nil {
			return nil, err
		}
	}
	return l, nil
}

func goagaginWait(l net.Listener) error {
	// Block the main goroutine awaiting signals.
	if _, err := goagain.Wait(l); nil != err {
		return err
	}

	// Do whatever's necessary to ensure a graceful exit like waiting for
	// goroutines to terminate or a channel to become closed.
	//
	// In this case, we'll simply stop listening and wait one second.
	return l.Close()
}

// A very rude server that says hello and then closes your connection.
func serve(e *echo.Echo, l net.Listener) {
	e.Run(standard.WithConfig(engine.Config{Address: l.Addr().String(), Listener: l}))
}
