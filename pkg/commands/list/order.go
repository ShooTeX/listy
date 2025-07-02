package list

import (
	"context"
	"fmt"

	"github.com/shootex/listy/internal/trakt"
	"github.com/urfave/cli/v3"
)

func Order() *cli.Command {
	var from string
	var destination string
	return &cli.Command{
		Name:  "order",
		Usage: "copy the order of items from one list to another",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "from",
				Destination: &from,
				Required:    true,
				Aliases:     []string{"f"},
			},
			&cli.StringFlag{
				Name:        "destination",
				Destination: &destination,
				Required:    true,
				Aliases:     []string{"to"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			traktApi, err := trakt.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create trakt client: %w", err)
			}

			if err := traktApi.CopyListOrder(from, destination); err != nil {
				return fmt.Errorf("failed to copy list order: %w", err)
			}

			return nil
		},
	}
}
