package commands

import (
	"context"
	"fmt"

	"github.com/shootex/listy/internal/auth"
	"github.com/urfave/cli/v3"
)

func configCmd() *cli.Command {
	var creds auth.Credentials
	return &cli.Command{
		Name:  "config",
		Usage: "setup trakt api credentials",
		Description: `This command allows you to set up your Trakt API credentials, 
which are required for authentication and accessing your Trakt account.
You will need to provide your client ID and client secret, 
which can be obtained from https://trakt.tv/oauth/applications.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "client-id",
				Required:    true,
				Aliases:     []string{"i"},
				Destination: &creds.ClientID,
			},
			&cli.StringFlag{
				Name:        "client-secret",
				Required:    true,
				Aliases:     []string{"s"},
				Destination: &creds.ClientSecret,
			},
		},
		Action: func(context.Context, *cli.Command) error {
			if err := auth.SaveCredentials(&creds); err != nil {
				return fmt.Errorf("failed to save credentials: %v", err)
			}
			return nil
		},
	}
}
