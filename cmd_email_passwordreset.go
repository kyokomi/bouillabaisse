package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/firebase"
)

func sendPasswordResetEmailCommand(c *cli.Context) error {
	fmt.Println()

	email := c.String("email")
	if email == "" {
		input, err := inputText("email")
		if err != nil {
			return err
		}
		email = input
	}

	return sendPasswordResetEmail(email)
}

func sendPasswordResetEmail(email string) error {
	fireClient := firebase.NewClient(
		firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
	)

	err := fireClient.Auth.SendPasswordResetEmail(email)
	if err != nil {
		return err
	}

	pp.Println("send ok")

	return nil
}
