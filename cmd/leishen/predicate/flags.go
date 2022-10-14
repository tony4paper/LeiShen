package predicate

import "github.com/urfave/cli/v2"

var (
	blockDBFlag = &cli.StringFlag{
		Name:     "blockdb",
		Usage:    "block database path",
		Aliases:  []string{"b"},
		Required: true,
	}

	receiptDBFlag = &cli.StringFlag{
		Name:     "receiptdb",
		Usage:    "receipt database path",
		Aliases:  []string{"r"},
		Required: true,
	}

	itxDBFlag = &cli.StringFlag{
		Name:     "itxdb",
		Usage:    "internal Transaction Database Path",
		Aliases:  []string{"i"},
		Required: true,
	}

	contractDBFlag = &cli.StringFlag{
		Name:     "contractdb",
		Usage:    "contract database path",
		Aliases:  []string{"c"},
		Required: true,
	}

	fltxDBFlag = &cli.StringFlag{
		Name:     "fltxdb",
		Usage:    "flash loan transaction database path",
		Aliases:  []string{"fl"},
		Required: true,
	}

	platformNameDBFlag = &cli.StringFlag{
		Name:     "pfndb",
		Usage:    "platform name database path",
		Aliases:  []string{"p"},
		Required: true,
	}

	txsFlag = &cli.StringSliceFlag{
		Name:    "txs",
		Usage:   "give the hash value of the transaction to be detected, which can be one or more",
		Aliases: []string{"t"},
	}

	txsFileFlag = &cli.StringFlag{
		Name:    "txsfile",
		Usage:   "given a file, where each row is a hash of transactions to be checked",
		Aliases: []string{"f"},
	}
)
