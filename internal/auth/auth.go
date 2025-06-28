package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"resty.dev/v3"
)

const (
	traktDeviceCodeURL = "https://api.trakt.tv/oauth/device/code"
	traktTokenURL      = "https://api.trakt.tv/oauth/device/token"
)

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

func StartDeviceAuthFlow(ctx context.Context) (*Token, error) {
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
		return nil, fmt.Errorf("failed to start device auth flow: %w", err)
	}

	deviceResp := resp.Result().(*DeviceCodeResponse)

	fmt.Printf("Please visit: %s/%s\n", deviceResp.VerificationURL, deviceResp.UserCode)
	fmt.Printf("Waiting for authorization...\n")

	ticker := time.NewTicker(time.Duration(deviceResp.Interval) * time.Second)
	defer ticker.Stop()

	timeout := time.After(time.Duration(deviceResp.ExpiresIn) * time.Second)

	for {
		select {
		case <-ticker.C:
			var tokenResp Token
			resp, err := http.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]string{
					"client_id":     clientId,
					"client_secret": clientSecret,
					"code":          deviceResp.DeviceCode,
				}).
				SetResult(&tokenResp).
				Post(traktTokenURL)
			if err != nil {
				return nil, fmt.Errorf("error during token exchange: %w", err)
			}

			if resp.IsError() {
				if resp.StatusCode() == 400 {
					continue // Authorization pending, keep waiting
				}
				return nil, fmt.Errorf("error during token exchange: %s", resp.String())
			}

			tokenResp.CreatedAt = UnixTime(time.Now())
			return &tokenResp, nil

		case <-timeout:
			return nil, errors.New("authorization timed out")
		}
	}
}
