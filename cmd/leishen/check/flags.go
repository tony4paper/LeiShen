package check

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

	fltxDBFlag = &cli.StringFlag{
		Name:     "fltxdb",
		Usage:    "flash loan transaction database path",
		Aliases:  []string{"fl"},
		Required: true,
	}

	stratFlag = &cli.Uint64Flag{
		Name:    "start",
		Usage:   "block height at which to start detection",
		Aliases: []string{"s"},
		Value:   0,
	}

	limitFlag = &cli.Uint64Flag{
		Name:    "limit",
		Usage:   "block height at which to stop detection",
		Aliases: []string{"l"},
		Value:   14500000,
	}

	forceFlag = &cli.BoolFlag{
		Name:    "force",
		Usage:   "whether to enable mandatory detection. After the detection is completed, the detected height range will be written. By default, the same range will not be repeatedly detected",
		Aliases: []string{"f"},
		Value:   false,
	}

	txsFlag = &cli.StringSliceFlag{
		Name:    "txs",
		Usage:   "give the hash value of the transaction to be detected, which can be one or more",
		Aliases: []string{"t"},
	}
)
