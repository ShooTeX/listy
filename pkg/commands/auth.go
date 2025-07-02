package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/shootex/listy/internal/auth"
	"github.com/urfave/cli/v3"
)

var AuthCmd = cli.Command{
	Name:  "auth",
	Usage: "Authenticate with Trakt",
	Action: func(ctx context.Context, c *cli.Command) error {
		token, err := auth.StartDeviceAuthFlow(ctx)
		if err != nil {
			return fmt.Errorf("failed to authenticate: %v", err)
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %v", err)
		}

		log.Println("Authentication successful!")
		return nil
	},
}
