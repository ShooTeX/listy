package list

import (
	"context"
	"fmt"

	"github.com/shootex/listy/internal/trakt"
	"github.com/urfave/cli/v3"
)

func Clean() *cli.Command {
	var list string
	var watched bool
	return &cli.Command{
		Name:  "clean",
		Usage: "clean up a list",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "watched",
				Destination: &watched,
				Aliases:     []string{"w"},
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "list",
				Destination: &list,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			traktApi, err := trakt.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create trakt client: %w", err)
			}

			cleanOptions := &trakt.CleanOptions{Watched: watched}

			if err := traktApi.Clean(list, cleanOptions); err != nil {
				return fmt.Errorf("failed to clean list %s: %w", list, err)
			}

			return nil
		},
	}
}
