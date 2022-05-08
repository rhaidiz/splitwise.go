package splitwise

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type Client interface {
	Users
	Groups
	Friends
	Expenses
}

type transport struct {
	underlyingTransport http.RoundTripper
	auth                AuthProvider
}

const (
	ServerAddress = "https://secure.splitwise.com"
)

// NewClient returns a new Client with the given AuthProvider
func NewClient(authProvider AuthProvider) Client {
	return &SClient{
		//AuthProvider: authProvider,
		BaseURL:    ServerAddress,
		HttpClient: http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, auth: authProvider}},
	}
}

func NewAuth0Client(conf *oauth2.Config, tok *oauth2.Token) Client {
	ctx := context.Background()
	client := conf.Client(ctx, tok)
	return &SClient{BaseURL: ServerAddress, HttpClient: *client, Conf: *conf}
}

type SClient struct {
	AuthProvider
	BaseURL    string
	HttpClient http.Client
	Conf       oauth2.Config
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, _ := t.auth.Auth()
	req.Header.Add("Authorization", "Bearer "+token)
	return t.underlyingTransport.RoundTrip(req)
}
