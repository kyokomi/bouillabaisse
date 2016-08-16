package main

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/store"
)

func removeLocalStoreCommand(c *cli.Context) error {
	fmt.Println()

	uid := c.String("uid")
	if uid == "" {
		input, err := inputText("uid")
		if err != nil {
			return err
		}
		uid = input
	}
	return removeLocalStore(uid)
}

func removeLocalStore(uid string) error {
	store.Stores.Remove(uid)
	pp.Printf("[%s] remove ok\n", uid)

	return nil
}
