package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shootex/listy/internal/auth"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cmd := &cli.Command{
		Name:                   "listy",
		Usage:                  "A simple CLI tool to manage Trakt lists",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
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
			},
			{
				Name:  "list",
				Usage: "Manage lists",
				Commands: []*cli.Command{
					addIntersectionToListCmd(),
					addDifferenceToListCmd(),
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
