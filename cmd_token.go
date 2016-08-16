package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/store"
)

func refreshTokenCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}
	return refreshToken(uid)
}

func refreshToken(uid string) error {
	authStore, ok := store.Stores.Data[uid]
	if !ok {
		return fmt.Errorf("uid = [%s] account is not found\n", uid)
	}

	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	token, err := fireClient.Token.ExchangeRefreshToken(authStore.RefreshToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ExchangeRefreshToken error [%s]\n", err.Error())
	}

	authStore.Token = token.AccessToken
	authStore.RefreshToken = token.RefreshToken
	authStore.ExpiresIn = token.ExpiresIn
	authStore.UpdateAt = time.Now()

	store.Stores.Add(authStore) // 上書き

	pp.Println(authStore)

	return nil
}
