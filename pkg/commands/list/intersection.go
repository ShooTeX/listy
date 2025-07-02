package list

import (
	"context"
	"fmt"

	"github.com/shootex/listy/internal/trakt"
	"github.com/urfave/cli/v3"
)

func Intersection() *cli.Command {
	var lists []string
	var destination string
	var clean bool
	return &cli.Command{
		Name: "intersection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "destination",
				Destination: &destination,
				Required:    true,
				Aliases:     []string{"to"},
			},
			&cli.BoolFlag{
				Name:        "clean",
				Usage:       "remove items from the destination list that are not in the intersection",
				Destination: &clean,
				Aliases:     []string{"c"},
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "lists",
				Destination: &lists,
				Min:         2,
				Max:         -1,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			traktApi, err := trakt.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create trakt client: %w", err)
			}

			if err := traktApi.AddIntersectToList(ctx, lists, destination, clean); err != nil {
				return fmt.Errorf("failed to add intersection to list: %w", err)
			}

			return nil
		},
	}
}
