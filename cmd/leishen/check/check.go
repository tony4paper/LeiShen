package check

import (
	"context"
	"fmt"
	"leishen/pkg/edb"
	"leishen/pkg/etypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

func checkTxs(blockDB, receiptDB, itxDB ethdb.Database, txsHash []common.Hash) error {
	for _, txHash := range txsHash {
		block, index := edb.ReadBlockByTxHash(blockDB, txHash)
		if block == nil || index == nil {
			return errTxNotFount(txHash)
		}

		tx := block.Transactions()[*index]

		receipts := edb.ReadReceipts(blockDB, receiptDB, block.NumberU64())
		if int(*index) >= len(receipts) {
			return errReceiptTxNotFount(block.NumberU64())
		}
		receipt := receipts[*index]

		itx := edb.ReadAllItxsWithNumber(itxDB, block.NumberU64(), txHash)

		fltx, err := CheckTransaction(tx.Hash(), receipt, itx)
		if err != nil {
			return err
		}

		fmt.Println(fltx)
	}
	return nil
}

const parallelism = 4

func checkRange(blockDB, receiptDB, itxDB, fltxDB ethdb.Database, start, limit uint64) error {
	save_start, _ := edb.ReadRange(fltxDB)
	if save_start != nil && *save_start > start {
		start = *save_start
	}

	bar := progressbar.Default(int64(limit-start), "check flash loan")

	grounp, ctx := errgroup.WithContext(context.Background())
	numberCh := make(chan uint64, parallelism)
	for i := 0; i < parallelism; i++ {
		grounp.Go(func() error {
			return checkFlashLoanInBlock(ctx, blockDB, receiptDB, itxDB, fltxDB, numberCh)
		})
	}

	grounp.Go(func() error {
		for number := start; number < limit; number++ {
			select {
			case <-ctx.Done():
				// 出错回滚
				edb.WriteRange(fltxDB, number-parallelism, number-parallelism)
				return nil

			case numberCh <- number:
				// 更新进度
				bar.Add(1)
				edb.WriteRange(fltxDB, number, number)
			}
		}

		// 通知结束
		close(numberCh)
		return nil
	})

	return grounp.Wait()
}

func checkFlashLoanInBlock(ctx context.Context, blockDB, receiptDB, itxDB, fltxDB ethdb.Database, numberCh chan uint64) error {
	var ok bool
	var number uint64
	for {
		select {
		case <-ctx.Done():
			return nil

		case number, ok = <-numberCh:
			if !ok {
				return nil
			}
		}

		block, receipts := edb.ReadBlockAndReceipts(blockDB, receiptDB, number)
		if block == nil {
			return errBlockNotFount(number)
		}

		if !BlockMayHaveFlashLoan(block) {
			continue
		}

		txs := block.Transactions()
		for i, tx := range txs {
			if !TransactionMayBeFlashLoan(receipts[i]) {
				continue
			}

			itx := edb.ReadAllItxsWithNumber(itxDB, block.NumberU64(), tx.Hash())
			fltx, err := CheckTransaction(tx.Hash(), receipts[i], itx)
			if err != nil {
				return err
			}

			if fltx != nil {
				edb.WriteFlashLoanTx(fltxDB, fltx)
			}
		}
	}
}

func BlockMayHaveFlashLoan(block *types.Block) bool {
	for _, feature := range AllFeatures {
		if feature.MayBeInBlock(block) {
			return true
		}
	}
	return false
}

func TransactionMayBeFlashLoan(receipt *types.Receipt) bool {
	for _, feature := range AllFeatures {
		if feature.MayBeInTransaction(receipt) {
			return true
		}
	}
	return false
}

var (
	bzx1 = &etypes.FlashLoan{
		TXHash:       common.HexToHash("0xb5c8bd9430b6cc87a0e2fe110ece6bf527fa4f170a4bc8cd032f768fc5219838"),
		Types2Loaner: map[string]common.Address{"bzx": common.HexToAddress("0x148426fdc4c8a51b96b4bed827907b5fa6491ad0")},
	}
	bzx2 = &etypes.FlashLoan{
		TXHash:       common.HexToHash("0x762881b07feb63c436dee38edd4ff1f7a74c33091e534af56c9f7d49b5ecac15"),
		Types2Loaner: map[string]common.Address{"bzx": common.HexToAddress("0xb8c6ad5fe7cb6cc72f2c4196dca11fbb272a8cbf")},
	}
)

func CheckTransaction(txHash common.Hash, receipt *types.Receipt, internalTxs []*etypes.InternalTransaction) (*etypes.FlashLoan, error) {
	if txHash == bzx1.TXHash {
		return bzx1, nil
	}

	if txHash == bzx2.TXHash {
		return bzx2, nil
	}

	m := map[string]common.Address{}
	for _, feature := range AllFeatures {
		if !feature.IsFlashLoan(receipt, internalTxs) {
			continue
		}

		name := feature.Name
		addr := getBorrowerAddress(name, receipt, internalTxs)
		if addr == nil {
			return nil, errBorrowerNotFound(txHash, name)
		}

		m[name] = *addr
	}

	if len(m) == 0 {
		return nil, nil
	}

	return &etypes.FlashLoan{
		TXHash:       txHash,
		Types2Loaner: m,
	}, nil
}
