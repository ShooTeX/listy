package commands

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/shootex/listy/internal/auth"
	"github.com/urfave/cli/v3"
	"github.com/zalando/go-keyring"
)

var authCmd = &cli.Command{
	Name:  "auth",
	Usage: "Authenticate with Trakt",
	Action: func(ctx context.Context, c *cli.Command) error {
		token, err := auth.StartDeviceAuthFlow(ctx)
		if err != nil {
			if errors.Is(err, keyring.ErrNotFound) {
				return fmt.Errorf("no client credentials found, please set them using 'listy config'")
			}
			return fmt.Errorf("failed to authenticate: %v", err)
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %v", err)
		}

		log.Println("Authentication successful!")
		return nil
	},
}
