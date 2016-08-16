package provider

// Config provider configuration
type Config struct {
	AuthSecretKey            string // gomniauth setup secretKey
	GitHubClientID           string
	GitHubSecretKey          string
	GoogleClientID           string
	GoogleSecretKey          string
	FacebookClientID         string
	FacebookSecretKey        string
	TwitterConsumerID        string
	TwitterConsumerSecretKey string
}
