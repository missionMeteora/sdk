package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PathDNA/ptk"
)

const (
	// Version is the current SDK version
	Version = "v0.5"

	// DefaultServer is the default api server used
	DefaultServer = "xxxhttps://meteora.us"
)

// common errors
var (
	ErrMissingUID = errors.New("missing user id")

	ErrMissingCID = errors.New("missing campaign id")

	ErrMissingSegID = errors.New("missing segment id")

	ErrMissingSegName = errors.New("missing segment name")

	ErrMissingProxSegName = errors.New("missing proximity segment name")

	ErrMissingLocations = errors.New("must specify at least one proximity location")

	ErrDateRange = errors.New("bad or missing date range")

	ErrInvalidName = errors.New("invalid name")

	ErrInvalidAgencyID = errors.New("invalid agency id")

	ErrInvalidEmail = errors.New("invalid email")

	ErrInvalidPassword = errors.New("password must be at least 8 characters")

	ErrInvalidLanding = errors.New("invalid landing url")

	ErrMissingAdImage = errors.New("missing ad image")

	ErrMissingAds = errors.New("must pass at least one ad")

	ErrInvalidAdSize = fmt.Errorf("invalid ad size, the accepted sizes are: %s", AllowedAdRectsString())

	ErrInvalidBudget = fmt.Errorf("invalid budget, must be between $%d and %d", CampaignMinimumBudget, CampaignMaximumBudget)

	ErrInvalidImpBudget = fmt.Errorf("invalid imp budget, must be less than %d", CampaignMaximumImpBudget)

	ErrInvalidSchedule = errors.New("invalid schedule provided")

	ErrRequestIsNil = errors.New("request is nil")
)

// New returns a new instance with the default server addr and given apiKey.
func New(apiKey string) *Client { return NewWithAddr(DefaultServer, apiKey) }

// NewWithAddr returns a new instance of Client with the given api addr and key.
// If testing locally on an unsecure https server, use https+insecure://localhost:8080.
func NewWithAddr(apiAddr string, apiKey string) *Client {
	u, err := url.Parse(apiAddr)
	if err != nil {
		log.Panicf("bad apiAddr %s: %v", apiAddr, err)
	}

	return newClient(u, apiKey)
}

func newClient(apiAddr *url.URL, apiKey string) *Client {
	c := &Client{
		c: ptk.HTTPClient{
			DefaultHeaders: http.Header{
				"X-SDK-VERSION": {Version},
				"X-APIKEY":      {apiKey},
			},
		},
		u: apiAddr, // less pointer derefs, and the Client itself is a pointer so we don't have to worry about copies.
	}

	if c.u.Scheme == "https+insecure" {
		c.c.AllowInsecureTLS(true)
		c.u.Scheme = "https"
	}

	c.u.Path = "/api/v1/"

	return c
}

// Client is a helper structure who holds onto an apiKey
type Client struct {
	u *url.URL
	c ptk.HTTPClient
}

// CurrentKey returns the API key used to initalize this client
func (c *Client) CurrentKey() string {
	return c.c.DefaultHeaders.Get("X-APIKEY")
}

// RawRequest is an alias for RawRequestCtx(context.Background(), method, endpoint, req, resp)
func (c *Client) RawRequest(method, endpoint string, req, resp interface{}) (err error) {
	return c.RawRequestCtx(context.Background(), method, endpoint, req, resp)
}

// RawRequestCtx allows a raw c.raw request to the given endpoint
func (c *Client) RawRequestCtx(ctx context.Context, method, endpoint string, req, resp interface{}) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var u *url.URL
	if u, err = c.getURL(endpoint); err != nil {
		return
	}

	roe := &respOrError{V: resp}

	ct := ""
	if method != "" {
		ct = "application/json"
	}

	return c.c.RequestCtx(ctx, method, ct, u.String(), req, roe.Handle)
}

func (c *Client) getURL(endpoint string) (u *url.URL, err error) {
	endpoint = strings.TrimPrefix(endpoint, "/")
	return c.u.Parse(endpoint)
}

func (c *Client) rawGet(ctx context.Context, endpoint string, resp interface{}) (err error) {
	return c.RawRequestCtx(ctx, "GET", endpoint, nil, resp)
}

func (c *Client) rawPut(ctx context.Context, endpoint string, req, resp interface{}) (err error) {
	return c.RawRequestCtx(ctx, "PUT", endpoint, req, resp)
}

func (c *Client) rawPost(ctx context.Context, endpoint string, req, resp interface{}) (err error) {
	return c.RawRequestCtx(ctx, "POST", endpoint, req, resp)
}

func (c *Client) rawDelete(ctx context.Context, endpoint string, resp interface{}) (err error) {
	return c.RawRequestCtx(ctx, "DELETE", endpoint, nil, resp)
}
