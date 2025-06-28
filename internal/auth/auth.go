package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"resty.dev/v3"
)

const (
	keyringService = "listy"
	keyringUser    = "trakt_oauth_token"

	traktDeviceCodeURL = "https://api.trakt.tv/oauth/device/code"
	traktTokenURL      = "https://api.trakt.tv/oauth/device/token"
)

func SaveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return keyring.Set(keyringService, keyringUser, string(data))
}

func LoadToken() (*oauth2.Token, error) {
	data, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func DeleteToken() error {
	return keyring.Delete(keyringService, keyringUser)
}

func getClientCredentials() (clientId, clientSecret string, err error) {
	clientId = os.Getenv("TRAKT_CLIENT_ID")
	clientSecret = os.Getenv("TRAKT_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		err = errors.New("TRAKT_CLIENT_ID or TRAKT_CLIENT_SECRET environment variables not set")
	}
	return
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func StartDeviceAuthFlow(ctx context.Context) (*oauth2.Token, error) {
	clientId, clientSecret, err := getClientCredentials()
	if err != nil {
		return nil, err
	}

	http := resty.New()

	resp, err := http.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"client_id": clientId}).
		SetResult(&DeviceCodeResponse{}).
		Post(traktDeviceCodeURL)
	if err != nil {
		return nil, err
	}

	deviceResp := resp.Result().(*DeviceCodeResponse)

	fmt.Printf("Please visit %s and enter the code: %s\n", deviceResp.VerificationURL, deviceResp.UserCode)
	fmt.Printf("Waiting for authorization...\n")

	ticker := time.NewTicker(time.Duration(deviceResp.Interval) * time.Second)
	defer ticker.Stop()

	timeout := time.After(time.Duration(deviceResp.ExpiresIn) * time.Second)

	for {
		select {
		case <-ticker.C:
			resp, err := http.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]string{
					"client_id":     clientId,
					"client_secret": clientSecret,
					"code":          deviceResp.DeviceCode,
				}).
				SetResult(&oauth2.Token{}).
				Post(traktTokenURL)
			if err != nil {
				return nil, err
			}

			if resp.IsError() {
				if resp.StatusCode() == 400 {
					continue // Authorization pending, keep waiting
				}
				return nil, fmt.Errorf("error during token exchange: %s", resp.String())
			}

			token := resp.Result().(*oauth2.Token)
			return token, nil

		case <-timeout:
			return nil, errors.New("authorization timed out")
		}
	}
}

func Config() (*oauth2.Config, error) {
	clientID, clientSecret, err := getClientCredentials()
	if err != nil {
		return nil, err
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: traktTokenURL,
		},
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		Scopes:      []string{"public", "lists"},
	}, nil
}
