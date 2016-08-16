package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/store"
)

func signUpWithSignInEmailPasswordCommand(c *cli.Context) error {
	fmt.Println()

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

	return signUpWithSignInEmailPassword(email, password)
}

func signUpWithSignInEmailPassword(email, password string) error {
	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	a, err := fireClient.Auth.SignInWithEmailAndPassword(email, password)
	if err != nil {
		pp.Println(err)
		// たぶん未登録なのでSignUpする
		a, err = fireClient.Auth.CreateUserWithEmailAndPassword(email, password)
		if err != nil {
			return err
		}
	}

	authStore := store.NewAuthStore(a)
	store.Stores.Add(authStore)

	pp.Println(authStore)

	return nil
}
