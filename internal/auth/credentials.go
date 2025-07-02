package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var credentialsKeyringKey = "trakt_credentials"

func SaveCredentials(creds *Credentials) error {
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}
	return keyring.Set(keyringService, credentialsKeyringKey, string(data))
}

func DeleteCredentials() error {
	return keyring.Delete(keyringService, credentialsKeyringKey)
}

func LoadCredentials() (*Credentials, error) {
	data, err := keyring.Get(keyringService, credentialsKeyringKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}
