package main

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/store"
)

func showLocalStoreAccountCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}
	return showLocalStoreAccount(uid)
}

func showLocalStoreAccount(uid string) error {
	a, ok := store.Stores.Data[uid]
	if !ok {
		return fmt.Errorf("uid = [%s] account is not found\n", uid)
	}
	pp.Println(a)
	return nil
}
