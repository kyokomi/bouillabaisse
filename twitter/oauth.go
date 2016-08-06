package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mrjones/oauth"
)

type OAuthClient struct {
	Consumer *oauth.Consumer

	callBackURL   string
	twitterTokens map[string]*oauth.RequestToken
}

func NewTwitterOAuth(cKey, cSecret, callBackURL string) *OAuthClient {
	consumer := oauth.NewConsumer(cKey, cSecret, oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	})
	anaconda.SetConsumerKey(cKey)
	anaconda.SetConsumerSecret(cSecret)
	return &OAuthClient{
		Consumer:      consumer,
		callBackURL:   callBackURL,
		twitterTokens: make(map[string]*oauth.RequestToken, 0),
	}
}

func (t OAuthClient) GetRequestTokenAndURL() (string, error) {
	token, requestURL, err := t.Consumer.GetRequestTokenAndUrl(t.callBackURL)
	if token != nil {
		t.addTwitterToken(token)
	}
	return requestURL, err
}

func (t *OAuthClient) GetAccessToken(tokenKey, verificationCode string) (*oauth.AccessToken, error) {
	return t.Consumer.AuthorizeToken(t.getTwitterToken(tokenKey), verificationCode)
}

func (t *OAuthClient) addTwitterToken(token *oauth.RequestToken) {
	t.twitterTokens[token.Token] = token
}

func (t *OAuthClient) getTwitterToken(token string) *oauth.RequestToken {
	return t.twitterTokens[token]
}
