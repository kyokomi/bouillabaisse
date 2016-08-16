package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/kyokomi/bouillabaisse/store"
)

func showLocalStoreAccountListCommand(_ *cli.Context) error {
	fmt.Println()

	return showLocalStoreAccountList()
}

func showLocalStoreAccountList() error {
	for _, a := range store.Stores.Data {
		fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%s\t%s\n",
			a.LocalID,
			a.ProviderID,
			a.DisplayName,
			a.ExpiresInText(),
			a.Email,
		)
	}
	return nil
}
