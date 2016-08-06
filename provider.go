package main

import "fmt"

const (
	callbackPath  = "/auth/callback"
	authLoginPath = "/auth/login"

	GitHubProvider   Provider = "github"
	TwitterProvider  Provider = "twitter"
	FacebookProvider Provider = "facebook"
	GoogleProvider   Provider = "google"
	UnknownProvider  Provider = "unkonwn"
)

type Provider string

var providers = map[string]Provider{
	GitHubProvider.Name():   GitHubProvider,
	TwitterProvider.Name():  TwitterProvider,
	FacebookProvider.Name(): FacebookProvider,
	GoogleProvider.Name():   GoogleProvider,
}

func NewProvider(providerName string) Provider {
	provider, ok := providers[providerName]
	if ok {
		return provider
	}
	return UnknownProvider
}

func (p Provider) Name() string {
	return string(p)
}

func (p Provider) CallbackURL(domain string) string {
	return fmt.Sprintf("%s%s/%s", domain, callbackPath, p.Name())
}

func (p Provider) SignInURL(domain string) string {
	return fmt.Sprintf("%s%s/%s", domain, authLoginPath, p.Name())
}
