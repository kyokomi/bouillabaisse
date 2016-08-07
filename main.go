package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kyokomi/bouillabaisse/firebase/provider"
)

var (
	configPath = flag.String("c", "./config.yaml",
		"configuration fila path yaml [default: ./config.yaml]")
	env = flag.String("e", "default",
		"env default")
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

	domain, err := Serve(*env, *configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// TODO: help command
	// TODO: show local Auth info
	// TODO: refreshToken refresh

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

	if _, err := inputWait("exit ? [y/n]"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("exit goodbye!")
}
