package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/kyokomi/bouillabaisse/server"
	"github.com/kyokomi/bouillabaisse/store"
)

func signInOAuthProviderCommand(c *cli.Context) error {
	fmt.Println()

	providerName := c.String("provider")
	if providerName == "" {
		input, err := inputText("provider [ twitter / facebook / github / google ]")
		if err != nil {
			return err
		}
		providerName = input
	}

	return signInOAuthProvider(providerName)
}

func signInOAuthProvider(providerName string) error {
	p := provider.New(providerName)
	if p == provider.UnknownProvider {
		return fmt.Errorf("Don't support provider [%s]\n", providerName)
	}

	return server.ProviderServeWithConfig(p, cfg, oauthCallbackListener)
}

func oauthCallbackListener(ctx echo.Context) error {
	linkProviderName := ctx.Param("provider")
	linkProvider := provider.New(linkProviderName)
	postBody, err := provider.BuildSignInPostBody(linkProvider, ctx.QueryParams())
	if err != nil {
		return errors.Wrapf(err, "%s BuildSignInPostBody error", linkProvider.Name())
	}

	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)
	a, err := fireClient.Auth.SignInWithOAuth(linkProvider, postBody)
	if err != nil {
		return errors.Wrapf(err, "%s SignInWithOAuth error", linkProvider.Name())
	}

	pp.Println(a)

	store.Stores.Add(store.NewAuthStore(a))

	return nil
}
