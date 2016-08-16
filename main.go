package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"gopkg.in/go-pp/pp.v2"

	"github.com/kyokomi/bouillabaisse/config"
	"github.com/kyokomi/bouillabaisse/firebase"
	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"github.com/kyokomi/bouillabaisse/server"
	"github.com/kyokomi/bouillabaisse/store"
)

var (
	configPath = flag.String("c", "./config.yaml",
		"configuration fila path yaml [default: ./config.yaml]")
	env = flag.String("e", "default",
		"env default")
)

func main() {
	flag.Parse()

	cfg := config.NewConfig(*env, *configPath)

	if err := store.Stores.Load(cfg.Local.AuthStoreDirPath, cfg.Local.AuthStoreFileName); err != nil {
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
			for _, a := range store.Stores.Data {
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
			if a, ok := store.Stores.Data[uid]; ok {
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
				firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
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

			emailStore := store.AuthStore{Auth: a, CreatedAt: time.Now(), UpdateAt: time.Now()}
			store.Stores.Add(emailStore)

			pp.Println(emailStore)

		case "link-email":
			uid := getInputSubCommand(input)
			if aStore, ok := store.Stores.Data[uid]; ok {
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
					firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
				)

				var linkAuth firebase.Auth
				linkAuth, err = fireClient.Auth.LinkAccountsAsyncWithEmailAndPassword(aStore.Auth, email, password)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				linkAuthStore := store.AuthStore{Auth: linkAuth, CreatedAt: aStore.CreatedAt, UpdateAt: time.Now()}
				store.Stores.Add(linkAuthStore)

				pp.Println(linkAuthStore)
			}

		case "link-oauth":
			uid := getInputSubCommand(input)

			if aStore, ok := store.Stores.Data[uid]; ok {

				providerName, err := inputText("link provider [ twitter / facebook / google / github ]")
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				p := provider.New(providerName)
				if p == provider.UnknownProvider {
					fmt.Fprintf(os.Stderr, "Don't support provider [%s]\n", providerName)
					os.Exit(1)
				}

				if err := server.ProviderServeWithConfig(p, cfg, func(ctx echo.Context) error {
					linkProviderName := ctx.Param("provider")
					linkProvider := provider.New(linkProviderName)
					postBody, err := provider.BuildSignInPostBody(linkProvider, ctx.QueryParams())
					if err != nil {
						return errors.Wrapf(err, "%s BuildSignInPostBody error", linkProvider.Name())
					}

					fireClient := firebase.NewClient(
						firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
					)

					var linkAuth firebase.Auth
					linkAuth, err = fireClient.Auth.LinkAccountsWithOAuth(aStore.Auth, postBody)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					linkAuthStore := store.AuthStore{Auth: linkAuth, CreatedAt: aStore.CreatedAt, UpdateAt: time.Now()}
					store.Stores.Add(linkAuthStore)

					pp.Println(linkAuthStore)

					return nil
				}); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

		case "anonymously":
			fireClient := firebase.NewClient(
				firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
			)

			a, err := fireClient.Auth.SignInAnonymously()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			aStore := store.AuthStore{Auth: a, CreatedAt: time.Now(), UpdateAt: time.Now()}
			store.Stores.Add(aStore)

			pp.Println(aStore)

		case "local-remove":
			uid := getInputSubCommand(input)

			store.Stores.Remove(uid)
			pp.Printf("[%s] remove ok\n", uid)

		case "token":
			uid := getInputSubCommand(input)

			a, ok := store.Stores.Data[uid]
			if !ok {
				fmt.Fprintf(os.Stderr, "Not found uid [%s]\n", uid)
			} else {
				fireClient := firebase.NewClient(
					firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
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

					store.Stores.Add(a) // 上書き
				}
			}

		case "provider":
			providerName := getInputSubCommand(input)

			p := provider.New(providerName)
			if p == provider.UnknownProvider {
				fmt.Fprintf(os.Stderr, "Don't support provider [%s]\n", providerName)
				os.Exit(1)
			}

			if err := server.ProviderServeWithConfig(p, cfg, func(ctx echo.Context) error {
				linkProviderName := ctx.Param("provider")
				linkProvider := provider.New(linkProviderName)
				postBody, err := provider.BuildSignInPostBody(linkProvider, ctx.QueryParams())
				if err != nil {
					return errors.Wrapf(err, "%s BuildSignInPostBody error", linkProvider.Name())
				}

				fireClient := firebase.NewClient(
					firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
				)
				auth, err := fireClient.Auth.SignInWithOAuth(linkProvider, postBody)
				if err != nil {
					return errors.Wrapf(err, "%s SignInWithOAuth error", linkProvider.Name())
				}

				pp.Println(auth)

				now := time.Now()
				a := store.AuthStore{Auth: auth, CreatedAt: now, UpdateAt: now}
				store.Stores.Add(a)
				return nil
			}); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

		case "email-verify":
			idToken := getInputSubCommand(input)
			fireClient := firebase.NewClient(
				firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
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
				firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
			)

			if err := fireClient.Auth.SendNewEmailAccept(store.Stores.Data[uid].Token,
				store.Stores.Data[uid].Email, nextEmail); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println("send ok")
			}
		case "password-reset":
			email := getInputSubCommand(input)
			fireClient := firebase.NewClient(
				firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
			)

			if err := fireClient.Auth.SendPasswordResetEmail(email); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println("send ok")
			}
		case "get-account":
			uid := getInputSubCommand(input)
			authStore := store.Stores.Data[uid]

			fireClient := firebase.NewClient(
				firebase.Config{APIKey: cfg.Server.FirebaseAPIKey}, &http.Transport{},
			)

			accountInfo, err := fireClient.Account.GetAccountInfo(authStore.Token)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				pp.Println(accountInfo)

				for _, u := range accountInfo.Users {
					if u.LocalID != authStore.LocalID {
						continue
					}

					authStore.DisplayName = u.DisplayName
					authStore.Email = u.Email
					authStore.PhotoURL = u.PhotoURL
					authStore.UpdateAt = time.Now()
					authStore.EmailVerified = u.EmailVerified

					store.Stores.Add(authStore)
				}
			}
		case "save":
			if err := store.Stores.Save(cfg.Local.AuthStoreDirPath, cfg.Local.AuthStoreFileName); err != nil {
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
