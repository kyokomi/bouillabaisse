package provider

type Config struct {
	AuthSecretKey            string // gomniauth setup secretKey
	GitHubClientID           string
	GitHubSecretKey          string
	GoogleClientID           string
	GoogleSecretKey          string
	FacebookID               string
	FacebookSecretKey        string
	TwitterConsumerID        string
	TwitterConsumerSecretKey string
}
