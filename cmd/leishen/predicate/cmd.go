package predicate

import (
	"bufio"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	PredicateCommand = &cli.Command{
		Action: predicate,
		Name:   "predicate",
		Usage:  "check if a transaction is a flash loan",
		Flags: []cli.Flag{
			blockDBFlag,
			receiptDBFlag,
			itxDBFlag,
			contractDBFlag,
			fltxDBFlag,
			platformNameDBFlag,
			txsFlag,
			txsFileFlag,
		},
	}
)

func predicate(ctx *cli.Context) error {
	blockDB, err := openDB(ctx.String(blockDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	receiptDB, err := openDB(ctx.String(receiptDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	itxDB, err := openDB(ctx.String(itxDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	fltxDB, err := openDB(ctx.String(fltxDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	platformNameDB, err := openDB(ctx.String(platformNameDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	contractDB, err := openDB(ctx.String(contractDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	if ctx.IsSet(txsFlag.Name) || ctx.IsSet(txsFileFlag.Name) {
		var txs []string
		if ctx.IsSet(txsFlag.Name) {
			txs = append(txs, ctx.StringSlice(txsFlag.Name)...)
		}

		if ctx.IsSet(txsFileFlag.Name) {
			contents, err := readHashFromFile(ctx.String(txsFileFlag.Name))
			if err != nil {
				return err
			}

			txs = append(txs, contents...)
		}

		if len(txs) > 0 {
			it := newTxHashIteratorFromStringSlice(txs)
			return anylyze(blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB, it)
		}
	}

	return anylyze(blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB, newTxHashIteratorFromFltxDB(fltxDB))
}

func readHashFromFile(file_name string) ([]string, error) {
	file, err := os.Open(file_name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var contents []string

	file_reader := bufio.NewScanner(file)
	for file_reader.Scan() {
		content := file_reader.Text()
		if content != "" {
			contents = append(contents, content)
		}
	}

	if err := file_reader.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}
