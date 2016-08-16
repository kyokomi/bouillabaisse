package main

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/store"
)

func saveLocalStoreCommand(_ *cli.Context) error {
	fmt.Println()

	return saveLocalStore()
}

func saveLocalStore() error {
	err := store.Stores.Save(cfg.Local.AuthStoreDirPath, cfg.Local.AuthStoreFileName)
	if err != nil {
		return err
	}

	pp.Println("save ok")

	return nil
}
