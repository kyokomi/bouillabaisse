package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"gopkg.in/go-pp/pp.v2"
)

var (
	configPath = flag.String("c", "./config.yaml",
		"configuration fila path yaml [default: ./config.yaml]")
	env = flag.String("e", "default",
		"env default")
	cmd = flag.String("cmd", "help", "execute command")
)

func inputWait(helpMessage string) (string, error) {
	fmt.Printf("%s > ", helpMessage)
	var input string
	if _, err := fmt.Scanf("%s", &input); err != nil {
		return "", err
	}
	return input, nil
}

func main() {
	flag.Parse()

	config := NewConfig(*env, *configPath)
	domain, err := ServeWithConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// TODO: help command
	// TODO: load AuthStore
	// TODO: show local Auth info
	// TODO: refreshToken refresh

	input := *cmd
	getInputCommand := func() string {
		return input
	}
	inputReset := func() {
		input = ""
	}
	inputSet := func(text string) {
		input = text
	}
	isEmptyCommand := func() bool {
		return getInputCommand() == ""
	}

	for true {
		if isEmptyCommand() {
			inputCmd, err := inputWait("[ exit / help / list / provider ]")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			inputSet(inputCmd)
		}

		switch getInputCommand() {
		case "help":
			fmt.Fprintln(os.Stdout, "help is comming soon!")
		case "list":
			a := AuthStore{}
			if err := a.Load(config.Local.AuthStoreDirPath, config.Local.AuthStoreFileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			pp.Println(a)
		case "provider":
			providerName, err := inputWait("provider [twitter/google/facebook/github/email]")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			p := provider.New(providerName)
			if p == provider.UnknownProvider {
				fmt.Fprintf(os.Stderr, "Don't support provider [%s]\n", providerName)
				os.Exit(1)
			}

			signInURL := provider.SignInURL(p, domain)
			fmt.Fprintln(os.Stdout, signInURL)
		case "exit":
			goto END
		}
		inputReset()
	}

END:
	fmt.Println("exit goodbye!")
}
