package commands

import (
	"github.com/shootex/listy/pkg/commands/list"
	"github.com/urfave/cli/v3"
)

var Commands = []*cli.Command{list.Cmd, &AuthCmd}
