package splitwise

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type Client interface {
	Users
	Groups
	Friends
	Expenses
	Oauth2Client
}

type Oauth2Client interface {
	GetOAuth2AuthorizeURL() string
	SetOAuth2Code(code string) error
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
	return &client{
		//AuthProvider: authProvider,
		baseURL: ServerAddress,
		client:  &http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, auth: authProvider}},
	}
}

func NewAuth0Client(conf *oauth2.Config) Client {
	return &client{conf: conf}
}

func (c client) GetOAuth2AuthorizeURL() string {
	url := c.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return url
}

func (c *client) SetOAuth2Code(code string) error {
	ctx := context.Background()
	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := c.conf.Exchange(ctx, code)
	if err != nil {
		return err
	}

	client := c.conf.Client(ctx, tok)
	c.baseURL = ServerAddress
	c.client = client

	return nil
}

type client struct {
	// AuthProvider
	baseURL string
	client  *http.Client
	conf    *oauth2.Config
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, _ := t.auth.Auth()
	req.Header.Add("Authorization", "Bearer "+token)
	return t.underlyingTransport.RoundTrip(req)
}
