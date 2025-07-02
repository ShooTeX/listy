package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
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

func DeleteToken() error {
	return keyring.Delete(keyringService, tokenKeyringKey)
}
