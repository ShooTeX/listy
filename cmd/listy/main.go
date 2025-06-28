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
		Name:  "listy",
		Usage: "A simple CLI tool to manage Trakt lists",
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
				Name: "test",
				Action: func(ctx context.Context, c *cli.Command) error {
					client, err := auth.NewClient(ctx, nil)
					if err != nil {
						return fmt.Errorf("failed to create client: %v", err)
					}
					defer client.Close()

					_, err = client.R().
						Get("/users/me/lists")
					if err != nil {
						log.Fatal(err)
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
