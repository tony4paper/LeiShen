package check

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

var (
	FlashLoanCheckCommand = &cli.Command{
		Action: check,
		Name:   "check",
		Usage:  "Check if the transaction is a flash loan transaction",
		Description: `
check checks a given transaction to see if it is a flash loan.
Currently, flash loan transactions on three platforms aave, dydx and uniswapV2 can be detected.
		`,
		Flags: []cli.Flag{
			blockDBFlag,
			receiptDBFlag,
			itxDBFlag,
			fltxDBFlag,
			stratFlag,
			limitFlag,
			forceFlag,
			txsFlag,
		},
	}
)

func check(c *cli.Context) error {
	blockDBPath := c.String(blockDBFlag.Name)
	blockDB, err := openDB(blockDBPath, 256, 64, true)
	if err != nil {
		return err
	}

	receiptDBPath := c.String(receiptDBFlag.Name)
	receiptDB, err := openDB(receiptDBPath, 256, 64, true)
	if err != nil {
		return err
	}

	itxDBPath := c.String(itxDBFlag.Name)
	itxDB, err := openDB(itxDBPath, 256, 64, true)
	if err != nil {
		return err
	}

	fltxDBPath := c.String(fltxDBFlag.Name)
	fltxDB, err := openDB(fltxDBPath, 256, 64, false)
	if err != nil {
		return err
	}

	if c.IsSet(txsFlag.Name) {
		txStringSlice := c.StringSlice(txsFlag.Name)
		txs := []common.Hash{}
		for _, tx := range txStringSlice {
			txs = append(txs, common.HexToHash(tx))
		}
		return checkTxs(blockDB, receiptDB, itxDB, txs)
	}

	start := c.Uint64(stratFlag.Name)
	limit := c.Uint64(limitFlag.Name)
	return checkRange(blockDB, receiptDB, itxDB, fltxDB, start, limit)
}
