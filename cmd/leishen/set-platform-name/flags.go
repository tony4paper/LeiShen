package setplatformname

import "github.com/urfave/cli/v2"

var (
	contractDBFlag = &cli.StringFlag{
		Name:     "contractdb",
		Usage:    "contract database path",
		Aliases:  []string{"c"},
		Required: true,
	}

	platformNameDBFlag = &cli.StringFlag{
		Name:     "pfndb",
		Usage:    "platform name database path",
		Aliases:  []string{"p"},
		Required: true,
	}

	csvFileFlag = &cli.StringFlag{
		Name:     "csv",
		Usage:    "csv file path containing the platform name, with two columns, the first is the account address and the second is the platform name",
		Aliases:  []string{"f"},
		Required: true,
	}
)
