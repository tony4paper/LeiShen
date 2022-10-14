package predicate

import (
	"context"
	"fmt"
	"leishen/cmd/leishen/check"
	"leishen/pkg/edb"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

const parallelism = 16

func anylyze(blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB ethdb.Database, it TxHashIterator) error {
	bar := progressbar.Default(-1, "Detecting flash loan attacks")

	grounp, ctx := errgroup.WithContext(context.Background())
	txHashCh := make(chan common.Hash, parallelism)
	for i := 0; i < parallelism; i++ {
		grounp.Go(func() error {
			return analyzeTransactionWithCh(ctx, blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB, txHashCh, nil)
		})
	}

	grounp.Go(func() error {
		for it.Next() {
			txHash := *it.TxHash()
			select {
			case <-ctx.Done():
				return nil

			case txHashCh <- txHash:
				bar.Add(1)
			}
		}

		close(txHashCh)
		return nil
	})

	return grounp.Wait()
}

func analyzeTransactionWithCh(ctx context.Context, blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB ethdb.Database, txHashCh chan common.Hash, w ResultWriter) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case txHash, ok := <-txHashCh:
			if !ok {
				return nil
			}

			err := analyzeTransaction(blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB, txHash, w)
			if err != nil {
				return err
			}
		}
	}
}

func analyzeTransaction(blockDB, receiptDB, itxDB, contractDB, fltxDB, platformNameDB ethdb.Database, txHash common.Hash, w ResultWriter) error {
	start := time.Now()
	var err error

	block, index := edb.ReadBlockByTxHash(blockDB, txHash)
	if block == nil || index == nil {
		return errTxNotFount(txHash)
	}

	tx := block.Transactions()[*index]
	message, err := tx.AsMessage(types.MakeSigner(params.MainnetChainConfig, block.Number()), nil)
	if err != nil {
		return errTxToMessage(txHash, err)
	}

	receipts := edb.ReadReceipts(blockDB, receiptDB, block.NumberU64())
	if int(*index) >= len(receipts) {
		return errReceiptTxNotFount(block.NumberU64())
	}
	receipt := receipts[*index]

	itx := edb.ReadAllItxsWithNumber(itxDB, block.NumberU64(), txHash)

	fltx := edb.ReadFlashLoanTx(fltxDB, txHash)
	if fltx == nil {
		fltx, err = check.CheckTransaction(txHash, receipt, itx)
		if err != nil {
			return err
		}

		if fltx == nil {
			return errTxNotFlashLoan(txHash)
		}
	}

	raw_ether_records := GetEtherRecords(itx)
	var ether_records []Record
	for _, record := range raw_ether_records {
		ether_records = append(ether_records, record)
	}
	ether_records = RemoveZeroRecord(ether_records)
	non_zero_itx_number := len(ether_records)

	raw_erc_records := GetErc20Records(receipt)
	erc_record_number := len(raw_erc_records)
	var erc_records []Record
	for _, record := range raw_erc_records {
		erc_records = append(erc_records, record)
	}

	erc_records = PreProcessRecord(contractDB, platformNameDB, fltx, erc_records)

	trades := SearchTrades(erc_records)
	p1 := SearchRaisingPrice(trades)
	p2 := SearchMultiRoundPriceFluctuation(trades)
	p3 := SearchCollateralizingLiquidity(erc_records)
	p4 := SearchBuyLowSellHigh(trades)

	fmt.Printf("%s, %d, %s, %s, %s, %s, %d, %d, %d, %d, %d, %d, %d\n", txHash, block.NumberU64(), message.From(), message.To(), time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"), strings.Join(fltx.Types(), "|"),
		non_zero_itx_number, erc_record_number, time.Since(start).Milliseconds(),
		len(p1), len(p2), len(p3), len(p4))

	return nil
}
