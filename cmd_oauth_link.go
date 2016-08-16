package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/kyokomi/bouillabaisse/server"
	"github.com/kyokomi/bouillabaisse/store"
)

func linkOAuthProviderCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}

	providerName := c.String("provider")
	if providerName == "" {
		input, err := inputText("provider [ twitter / facebook / github / google ]")
		if err != nil {
			return err
		}
		providerName = input
	}

	return linkOAuthProvider(uid, providerName)
}

func linkOAuthProvider(uid, providerName string) error {
	p := provider.New(providerName)
	if p == provider.UnknownProvider {
		return fmt.Errorf("Don't support provider [%s]\n", providerName)
	}

	authStore, ok := store.Stores.Data[uid]
	if !ok {
		return fmt.Errorf("uid = [%s] account is not found\n", uid)
	}

	if err := server.ProviderServeWithConfig(p, cfg, linkOAuthCallbackListenerFunc(authStore)); err != nil {
		return err
	}

	return server.ProviderServeWithConfig(p, cfg, oauthCallbackListener)
}

func linkOAuthCallbackListenerFunc(authStore store.AuthStore) server.AuthCallbackListenerFunc {
	return func(ctx echo.Context) error {
		linkProviderName := ctx.Param("provider")
		linkProvider := provider.New(linkProviderName)
		postBody, err := provider.BuildSignInPostBody(linkProvider, ctx.QueryParams())
		if err != nil {
			return errors.Wrapf(err, "%s BuildSignInPostBody error", linkProvider.Name())
		}

		fireClient := firebase.NewClient(
			firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
		)

		var linkAuth firebase.Auth
		linkAuth, err = fireClient.Auth.LinkAccountsWithOAuth(authStore.Auth, postBody)
		if err != nil {
			return err
		}
		linkAuthStore := store.AuthStore{Auth: linkAuth, CreatedAt: authStore.CreatedAt, UpdateAt: time.Now()}
		store.Stores.Add(linkAuthStore)

		pp.Println(linkAuthStore)

		return nil
	}
}
