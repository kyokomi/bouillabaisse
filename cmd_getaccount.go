package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/store"
)

func getAccountCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}
	return getAccount(uid)
}

func getAccount(uid string) error {
	authStore, ok := store.Stores.Data[uid]
	if !ok {
		return fmt.Errorf("uid = [%s] account is not found\n", uid)
	}

	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	accountInfo, err := fireClient.Account.GetAccountInfo(authStore.Token)
	if err != nil {
		return err
	}

	for _, u := range accountInfo.Users {
		if u.LocalID != authStore.LocalID {
			continue
		}

		authStore.DisplayName = u.DisplayName
		authStore.Email = u.Email
		authStore.PhotoURL = u.PhotoURL
		authStore.UpdateAt = time.Now()
		authStore.EmailVerified = u.EmailVerified

		store.Stores.Add(authStore)

		pp.Println(authStore)
	}

	return nil
}
