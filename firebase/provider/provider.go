package provider

type Provider string

const (
	GitHubProvider   Provider = "github"
	TwitterProvider  Provider = "twitter"
	FacebookProvider Provider = "facebook"
	GoogleProvider   Provider = "google"
	UnknownProvider  Provider = "unkonwn"
)

var providerMaps = map[string]Provider{
	GitHubProvider.Name():   GitHubProvider,
	TwitterProvider.Name():  TwitterProvider,
	FacebookProvider.Name(): FacebookProvider,
	GoogleProvider.Name():   GoogleProvider,
}

var providerIDMaps = map[Provider]string{
	FacebookProvider: "facebook.com",
	GoogleProvider:   "google.com",
	GitHubProvider:   "github.com",
	TwitterProvider:  "twitter.com",
}

func New(providerName string) Provider {
	provider, ok := providerMaps[providerName]
	if ok {
		return provider
	}
	return UnknownProvider
}

func (p Provider) Name() string {
	return string(p)
}

func (p Provider) id() string {
	return providerIDMaps[p]
}
