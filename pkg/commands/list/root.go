package list

import "github.com/urfave/cli/v3"

var Cmd = &cli.Command{
	Name:  "list",
	Usage: "Manage lists",
	Commands: []*cli.Command{
		Intersection(),
		Difference(),
		Order(),
		Clean(),
	},
}
