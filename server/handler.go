package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"github.com/kyokomi/bouillabaisse/config"
	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
)

type serverContext struct {
	fireClient *firebase.Client
	config     config.Config
}

func (*serverContext) loginHandler(ctx echo.Context) error {
	providerName := ctx.Param("provider")

	loginURL, err := provider.GetBeginAuthURL(provider.New(providerName))
	if err != nil {
		return errors.Wrapf(err, "GetBeginAuthURLの呼び出し中にエラーが発生しました: %s", providerName)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (s *serverContext) callbackHandler(ctx echo.Context, callbackFunc AuthCallbackFunc) error {
	providerName := ctx.Param("provider")

	p := provider.New(providerName)

	if err := callbackFunc(ctx); err != nil {
		return errors.Wrapf(err, "%s callbackFunc error", p.Name())
	}

	return ctx.JSON(http.StatusOK, providerName)
}
