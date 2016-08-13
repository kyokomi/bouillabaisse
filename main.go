package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"time"

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

	config := newConfig(*env, *configPath)
	domain, err := serveWithConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := stores.Load(config.Local.AuthStoreDirPath, config.Local.AuthStoreFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
			for _, a := range stores.stores {
				fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%s\t%s\n",
					a.LocalID,
					a.ProviderID,
					a.DisplayName,
					a.ExpiresInText(),
					a.Email,
				)
			}

		case "show":
			uid := getInputSubCommand(input)
			if a, ok := stores.stores[uid]; ok {
				pp.Println(a)
			}
		case "email":
			email, err := inputText("email")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			password, err := inputText("password")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fireClient := firebase.NewClient(
				firebase.Config{APIKey: config.Server.FirebaseAPIKey}, &http.Transport{},
			)

			var a firebase.Auth
			a, err = fireClient.Auth.SignInWithEmailAndPassword(email, password)
			if err != nil {
				a, err = fireClient.Auth.CreateUserWithEmailAndPassword(email, password)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
			stores.Add(authStore{Auth: a, CreatedAt: time.Now(), UpdateAt: time.Now()})

			pp.Println(a)

		case "token":
			uid := getInputSubCommand(input)

			a, ok := stores.stores[uid]
			if !ok {
				fmt.Fprintf(os.Stderr, "Not found uid [%s]\n", uid)
			} else {
				fireClient := firebase.NewClient(
					firebase.Config{APIKey: config.Server.FirebaseAPIKey}, &http.Transport{},
				)
				token, err := fireClient.Token.ExchangeRefreshToken(a.RefreshToken)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ExchangeRefreshToken error [%s]\n", err.Error())
				} else {
					pp.Println(token)

					a.Token = token.AccessToken
					a.RefreshToken = token.RefreshToken
					a.ExpiresIn = token.ExpiresIn
					a.UpdateAt = time.Now()

					stores.Add(a) // 上書き
				}
			}

		case "provider":
			providerName := getInputSubCommand(input)

			p := provider.New(providerName)
			if p == provider.UnknownProvider {
				fmt.Fprintf(os.Stderr, "Don't support provider [%s]\n", providerName)
				os.Exit(1)
			}

			signInURL := provider.BuildSignInURL(p, domain)
			fmt.Fprintln(os.Stdout, signInURL)
		case "email-verify":
			idToken := getInputSubCommand(input)
			fireClient := firebase.NewClient(
				firebase.Config{APIKey: config.Server.FirebaseAPIKey}, &http.Transport{},
			)

			if err := fireClient.Auth.SendEmailVerify(idToken); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println("send ok")
			}
		case "new-email":
			uid := getInputSubCommand(input)

			nextEmail, err := inputText("next email")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fireClient := firebase.NewClient(
				firebase.Config{APIKey: config.Server.FirebaseAPIKey}, &http.Transport{},
			)

			if err := fireClient.Auth.SendNewEmailAccept(stores.stores[uid].Token,
				stores.stores[uid].Email, nextEmail); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println("send ok")
			}
		case "pasword-reset":
			email := getInputSubCommand(input)
			fireClient := firebase.NewClient(
				firebase.Config{APIKey: config.Server.FirebaseAPIKey}, &http.Transport{},
			)

			if err := fireClient.Auth.SendPasswordResetEmail(email); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println("send ok")
			}
		case "save":
			if err := stores.Save(config.Local.AuthStoreDirPath, config.Local.AuthStoreFileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println("save ok")
			}
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
	fmt.Printf("\n %s > ", message)
	if _, err := fmt.Scan(&text); err != nil {
		return "", err
	}
	return text, nil
}

func printCommandInput() {
	fmt.Print("\n [ exit / help / show <uid> / list / email ] \n [ provider / token <uid>/ email-verify <token> / new-email <uid> / pasword-reset <email> ]\n > ")
}
