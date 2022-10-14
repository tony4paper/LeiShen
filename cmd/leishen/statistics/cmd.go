package statistics

import (
	"fmt"
	"leishen/pkg/edb"
	"leishen/pkg/etypes"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

var (
	StatisticsCommand = &cli.Command{
		Action: statistics,
		Name:   "statistics",
		Usage:  "Statistics flash loan transaction data",
		Flags: []cli.Flag{
			blockDBFlag,
			fltxDBFlag,
		},
	}
)

type TxIndex struct {
	fltx   *etypes.FlashLoan
	header *types.Header
	index  uint64
}

func statistics(c *cli.Context) error {
	blockDB, err := openDB(c.String(blockDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}

	fltxDB, err := openDB(c.String(fltxDBFlag.Name), 256, 64, true)
	if err != nil {
		return err
	}
	defer fltxDB.Close()

	fltxs := edb.ReadAllFlashLoan(fltxDB)

	blockMap := map[uint64]*types.Header{}

	txs := []TxIndex{}
	bar := progressbar.Default(int64(len(fltxs)), "read block")
	for _, fltx := range fltxs {
		block, index, err := readNumbers(blockDB, fltx.TXHash)
		if err != nil {
			return err
		}

		if saved := blockMap[block.NumberU64()]; saved != nil {
			txs = append(txs, TxIndex{fltx, saved, index})
			bar.Add(1)
			continue
		}

		blockMap[block.NumberU64()] = block.Header()
		txs = append(txs, TxIndex{fltx, block.Header(), index})
		bar.Add(1)
	}

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].header.Number.Uint64() < txs[j].header.Number.Uint64() || (txs[i].header.Number.Uint64() == txs[j].header.Number.Uint64() && txs[i].index < txs[j].index)
	})

	fmt.Println("tx_hash, number, time, index, platform")
	for _, tx := range txs {
		fmt.Printf("%s, %d, %s, %d, %s\n", tx.fltx.TXHash, tx.header.Number.Uint64(), time.Unix(int64(tx.header.Time), 0).Format("2006-01-02 15:04:05"), tx.index, strings.Join(tx.fltx.Types(), "|"))
	}

	return nil
}

func readNumbers(db ethdb.Database, txHash common.Hash) (*types.Block, uint64, error) {
	block, index := edb.ReadBlockByTxHash(db, txHash)
	if block == nil || index == nil {
		return nil, 0, errTxNotFount(txHash)
	}

	return block, *index, nil
}
