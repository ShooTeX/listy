package main

import (
	"context"
	"fmt"

	"github.com/shootex/listy/internal/trakt"
	"github.com/urfave/cli/v3"
)

func addIntersectionToListCmd() *cli.Command {
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

			if err := traktApi.AddIntersectToList(lists, destination, clean); err != nil {
				return fmt.Errorf("failed to add intersection to list: %w", err)
			}

			return nil
		},
	}
}

func addDifferenceToListCmd() *cli.Command {
	var lists []string
	var destination string
	var clean bool
	return &cli.Command{
		Name: "difference",
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

			if err := traktApi.AddDifferenceToList(ctx, lists, destination, clean); err != nil {
				return fmt.Errorf("failed to add intersection to list: %w", err)
			}

			return nil
		},
	}
}

func copyListOrderCmd() *cli.Command {
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

func cleanCmd() *cli.Command {
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
