package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"gopkg.in/go-pp/pp.v2"
)

var (
	configPath = flag.String("c", "./config.yaml",
		"configuration fila path yaml [default: ./config.yaml]")
	env = flag.String("e", "default",
		"env default")
)

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

	getInputCommand := func(input string) string {
		commands := strings.Fields(input)
		if len(commands) >= 1 {
			return commands[0]
		}
		return ""
	}
	getInputSubCommand := func(input string) string {
		commands := strings.Fields(input)
		if len(commands) >= 2 {
			return commands[1]
		}
		return ""
	}

	printCommandInput()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		switch getInputCommand(input) {
		case "help":
			fmt.Println("[ exit / help / list / provider / show / token ]")
		case "list":
			a := AuthStore{}
			if err := a.Load(config.Local.AuthStoreDirPath, config.Local.AuthStoreFileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%s\n",
				a.LocalID,
				a.ProviderID,
				a.DisplayName,
				a.ExpiresInText(),
			)

		case "show":
			a := AuthStore{}
			if err := a.Load(config.Local.AuthStoreDirPath, config.Local.AuthStoreFileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if a.LocalID == getInputSubCommand(input) {
				pp.Println(a)
			}
		case "token":
			uid := getInputSubCommand(input)

			a := AuthStore{}
			if err := a.Load(config.Local.AuthStoreDirPath, config.Local.AuthStoreFileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if a.LocalID != uid {
				fmt.Fprintf(os.Stderr, "Not found uid [%s]\n", uid)
			} else {
				fireClient := firebase.NewClient(
					firebase.Config{ApiKey: config.Server.FirebaseApiKey}, &http.Transport{},
				)
				token, err := fireClient.Token.ExchangeRefreshToken(a.RefreshToken)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ExchangeRefreshToken error [%s]\n", err.Error())
				} else {
					pp.Println(token)
				}
			}

		case "provider":
			providerName, err := inputText("provider [twitter/google/facebook/github/email]")
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
		default:
			fmt.Fprintf(os.Stdout, "%s command not found\n", getInputCommand(input))
		}

		printCommandInput()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

END:
	fmt.Println("exit goodbye!")
}

func inputText(message string) (string, error) {
	var text string
	fmt.Printf("\n %s > \n", message)
	if _, err := fmt.Scan(&text); err != nil {
		return "", err
	}
	return text, nil
}

func printCommandInput() {
	fmt.Print("\n [ exit / help / list / provider ] > ")
}
