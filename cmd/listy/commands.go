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
	return &cli.Command{
		Name: "intersection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "destination",
				Destination: &destination,
				Required:    true,
				Aliases:     []string{"to"},
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
			trakt, err := trakt.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create trakt client: %w", err)
			}

			if err := trakt.AddIntersectToList(lists, destination); err != nil {
				return fmt.Errorf("failed to add intersection to list: %w", err)
			}

			return nil
		},
	}
}

func addDifferenceToListCmd() *cli.Command {
	var lists []string
	var destination string
	return &cli.Command{
		Name: "difference",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "destination",
				Destination: &destination,
				Required:    true,
				Aliases:     []string{"to"},
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
			trakt, err := trakt.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create trakt client: %w", err)
			}

			if err := trakt.AddDifferenceToList(lists, destination); err != nil {
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
			trakt, err := trakt.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create trakt client: %w", err)
			}

			if err := trakt.CopyListOrder(from, destination); err != nil {
				return fmt.Errorf("failed to copy list order: %w", err)
			}

			return nil
		},
	}
}
