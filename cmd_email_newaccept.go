package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/store"
)

func sendNewEmailAcceptCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}
	newEmail := c.String("email")
	if newEmail == "" {
		input, err := inputText("new email")
		if err != nil {
			return err
		}
		newEmail = input
	}

	return sendNewEmailAccept(uid, newEmail)
}

func sendNewEmailAccept(uid, newEmail string) error {
	authStore, ok := store.Stores.Data[uid]
	if !ok {
		return fmt.Errorf("uid = [%s] account is not found\n", uid)
	}

	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	err := fireClient.Auth.SendNewEmailAccept(authStore.Token, authStore.Email, newEmail)
	if err != nil {
		return err
	}

	pp.Println("send ok")

	return nil
}
