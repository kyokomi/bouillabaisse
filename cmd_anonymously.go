package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/store"
)

func signUpAnonymouslyCommand(_ *cli.Context) error {
	fmt.Println()

	return signUpAnonymously()
}

func signUpAnonymously() error {
	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	a, err := fireClient.Auth.SignInAnonymously()
	if err != nil {
		return err
	}
	authStore := store.NewAuthStore(a)
	store.Stores.Add(authStore)

	pp.Println(authStore)

	return nil
}
