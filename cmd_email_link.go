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

func linkEmailCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}

	email := c.String("email")
	if email == "" {
		input, err := inputText("email")
		if err != nil {
			return err
		}
		email = input
	}

	password := c.String("password")
	if password == "" {
		input, err := inputText("password")
		if err != nil {
			return err
		}
		password = input
	}

	return linkEmail(uid, email, password)
}

func linkEmail(uid, email, password string) error {
	authStore, ok := store.Stores.Data[uid]
	if !ok {
		return fmt.Errorf("uid = [%s] account is not found\n", uid)
	}

	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	linkAuth, err := fireClient.Auth.LinkAccountsAsyncWithEmailAndPassword(authStore.Auth, email, password)
	if err != nil {
		return err
	}
	linkAuthStore := store.AuthStore{Auth: linkAuth, CreatedAt: authStore.CreatedAt, UpdateAt: time.Now()}
	store.Stores.Add(linkAuthStore)

	pp.Println(linkAuthStore)

	return nil
}
