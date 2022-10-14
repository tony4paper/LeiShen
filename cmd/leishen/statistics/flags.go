package statistics

import "github.com/urfave/cli/v2"

var (
	blockDBFlag = &cli.StringFlag{
		Name:     "blockdb",
		Usage:    "block database path",
		Aliases:  []string{"b"},
		Required: true,
	}

	fltxDBFlag = &cli.StringFlag{
		Name:     "fltxdb",
		Usage:    "flash loan transaction database path",
		Aliases:  []string{"fl"},
		Required: true,
	}
)
