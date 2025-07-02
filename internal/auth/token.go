package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"
	"resty.dev/v3"
)

var tokenKeyringKey = "trakt_oauth_token"

type Token struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        string   `json:"scope"`
	CreatedAt    UnixTime `json:"created_at"`
	TokenType    string   `json:"token_type"`
}

func (t *Token) IsExpired() bool {
	return t.CreatedAt.Time().Add(time.Duration(t.ExpiresIn) * time.Second).Before(time.Now())
}

func SaveToken(token *Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	return keyring.Set(keyringService, tokenKeyringKey, string(data))
}

func LoadToken() (*Token, error) {
	data, err := keyring.Get(keyringService, tokenKeyringKey)
	if err != nil {
		return nil, err
	}

	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

func RefreshToken(token *Token) (*Token, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	http := resty.New()

	var newToken Token
	_, err = http.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"client_id":     creds.ClientID,
			"client_secret": creds.ClientSecret,
			"refresh_token": token.RefreshToken,
			"grant_type":    "refresh_token",
			"redirect_uri":  "urn:ietf:wg:oauth:2.0:oob",
		}).
		SetResult(&newToken).
		Post(traktTokenURL)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	if err := SaveToken(&newToken); err != nil {
		return nil, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return &newToken, nil
}

func DeleteToken() error {
	return keyring.Delete(keyringService, tokenKeyringKey)
}
