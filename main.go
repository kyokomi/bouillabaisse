package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/kyokomi/bouillabaisse/config"
	"github.com/kyokomi/bouillabaisse/store"
)

var cfg config.Config

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./config.yaml",
			Usage: "configuration fila path yaml",
		},
		cli.StringFlag{
			Name:  "env, e",
			Value: "default",
			Usage: "Cli application environment",
		},
	}

	app.Before = func(c *cli.Context) error {
		configPath := c.String("config")
		env := c.String("env")
		cfg = config.NewConfig(env, configPath)
		if err := store.Stores.Load(cfg.Local.AuthStoreDirPath, cfg.Local.AuthStoreFileName); err != nil {
			return err
		}
		return nil
	}
	app.Action = dialogueMode
	app.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "Shows a list of firebase accounts at local store",
			Action: showLocalStoreAccountListCommand,
		},
		{
			Name:   "show",
			Usage:  "Show a firebase account detail at local store",
			Action: showLocalStoreAccountCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
			},
		},
		{
			Name:   "email",
			Usage:  "Firebase Auth SignUp with SignIn email and password",
			Action: signUpWithSignInEmailPasswordCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "email,e",
					Usage: "email address",
				},
				cli.StringFlag{
					Name:  "password,p",
					Usage: "password",
				},
			},
		},
		{
			Name:   "anonymously",
			Usage:  "Firebase Auth SignUp Anonymously",
			Action: signUpAnonymouslyCommand,
		},
		{
			Name:   "oauth",
			Usage:  "Firebase Auth SignIn OAuth Account",
			Action: signInOAuthProviderCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "provider,p",
					Usage: "provider name [ twitter / facebook / github / google ]",
				},
			},
		},
		{
			Name:   "link-oauth",
			Usage:  "Firebase Auth Link OAuth Account",
			Action: linkOAuthProviderCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
				cli.StringFlag{
					Name:  "provider,p",
					Usage: "provider name [ twitter / facebook / github / google ]",
				},
			},
		},
		{
			Name:   "link-email",
			Usage:  "Firebase Auth Link Email Account",
			Action: linkEmailCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
				cli.StringFlag{
					Name:  "email,e",
					Usage: "link email address",
				},
				cli.StringFlag{
					Name:  "password,p",
					Usage: "password",
				},
			},
		},
		{
			Name:   "refresh-token",
			Usage:  "Refresh token for Firebase Auth Account",
			Action: refreshTokenCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
			},
		},
		{
			Name:   "get-account",
			Usage:  "Get Firebase Account info at Firebase Auth ",
			Action: getAccountCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
			},
		},
		{
			Name:   "new-email",
			Usage:  "Send New Email address accept for Email Account",
			Action: sendNewEmailAcceptCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
				cli.StringFlag{
					Name:  "email,e",
					Usage: "new email address",
				},
			},
		},
		{
			Name:   "email-verify",
			Usage:  "Send Email address Verify for Email Account",
			Action: sendEmailVerifyCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
			},
		},
		{
			Name:   "password-reset",
			Usage:  "Send Password Reset for Email account",
			Action: sendPasswordResetEmailCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "email,e",
					Usage: "email address",
				},
			},
		},
		{
			Name:   "local-remove",
			Usage:  "Removed a uid at local store",
			Action: removeLocalStoreCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uid",
					Usage: "firebase auth user id",
				},
			},
		},
		{
			Name:   "save",
			Usage:  "Saved accounts at local store",
			Action: saveLocalStoreCommand,
		},
	}
	app.Run(os.Args)
}

func dialogueMode(c *cli.Context) error {
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

	// show help
	if err := c.App.Command("help").Run(c); err != nil {
		return err
	}

	printCommandInput()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		command := getInputCommand(input)
		subCommand := getInputSubCommand(input)

		cm := c.App.Command(command)
		if cm != nil {
			if subCommand == "--help" || subCommand == "help" || subCommand == "h" {
				cli.ShowCommandHelp(c, command)
			} else {
				if err := cm.Run(c); err != nil {
					return err
				}
			}
		} else {
			switch command {
			case "exit":
				fmt.Println("exit goodbye!")
				return nil
			default:
				fmt.Fprintf(os.Stdout, "%s command not found\n", getInputCommand(input))
			}
		}
		printCommandInput()
	}
	return scanner.Err()
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
	fmt.Print("\n > ")
}
