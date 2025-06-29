package auth

import (
	"context"
	"fmt"

	"github.com/shootex/listy/internal/version"
	"resty.dev/v3"
)

// client is the authenticated HTTP client for Trakt API
type Client struct {
	resty    *resty.Client
	token    *Token
	ctx      context.Context
	onUpdate func(newToken *Token) error
}

// NewClient creates a new Trakt client with automatic token management
func NewClient(ctx context.Context, onUpdate func(newToken *Token) error) (*resty.Client, error) {
	clientId, _, err := getClientCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get client credentials: %w", err)
	}
	token, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("no token found, please authenticate first: %w", err)
	}

	c := &Client{
		resty:    resty.New(),
		token:    token,
		ctx:      ctx,
		onUpdate: onUpdate,
	}

	c.resty.
		SetBaseURL("https://api.trakt.tv").
		SetHeader("Content-Type", "application/json").
		SetHeader("trakt-api-version", "2").
		SetHeader("trakt-api-key", clientId).
		SetHeader("User-Agent", fmt.Sprintf("%s/%s", version.Name, version.Version)).
		AddRequestMiddleware(c.authMiddleware).
		AddResponseMiddleware(errorResponseMiddleware)

	return c.resty, nil
}

func (c *Client) authMiddleware(client *resty.Client, req *resty.Request) error {
	if c.token.IsExpired() {
		newToken, err := RefreshToken(c.ctx, c.token)
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		c.token = newToken
		if c.onUpdate != nil {
			if err := c.onUpdate(newToken); err != nil {
				return err
			}
		}
	}

	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	return nil
}

func errorResponseMiddleware(client *resty.Client, resp *resty.Response) error {
	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.Status())
	}
	return nil
}
