package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"
	"resty.dev/v3"
)

const (
	keyringService = "listy"
	keyringUser    = "trakt_oauth_token"
)

type Token struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        string   `json:"scope"`
	CreatedAt    UnixTime `json:"created_at"`
	TokenType    string   `json:"token_type"`
}

func (t *Token) IsExpired() bool {
	expiry := t.CreatedAt.Time().Add(time.Duration(t.ExpiresIn) * time.Second)
	return time.Now().After(expiry.Add(-1 * time.Minute))
}

func SaveToken(token *Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	return keyring.Set(keyringService, keyringUser, string(data))
}

func LoadToken() (*Token, error) {
	data, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil, err
	}

	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

func DeleteToken() error {
	return keyring.Delete(keyringService, keyringUser)
}

func RefreshToken(ctx context.Context, token *Token) (*Token, error) {
	clientId, clientSecret, err := getClientCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get client credentials: %w", err)
	}

	http := resty.New()

	http.SetDebug(true)

	var newToken Token
	resp, err := http.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"refresh_token": token.RefreshToken,
			"grant_type":    "refresh_token",
		}).
		SetResult(&newToken).
		Post(traktTokenURL)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to refresh token: %s", resp.String())
	}

	newToken.CreatedAt = UnixTime(time.Now())
	return &newToken, nil
}
