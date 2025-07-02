package main

import (
	"context"
	"log"
	"os"

	"github.com/shootex/listy/pkg/commands"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cmd := &cli.Command{
		Name:                   "listy",
		Usage:                  "A simple CLI tool to manage Trakt lists",
		UseShortOptionHandling: true,
		Commands:               commands.Commands,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
