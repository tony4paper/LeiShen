package edb

import (
	"fmt"
	"leishen/pkg/etypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

const (
	flashLoanPrefix       = byte(0x01)
	flashLoanPrefixLength = 33
)

func KeyFlashLoan(hash common.Hash) []byte {
	return append([]byte{flashLoanPrefix}, hash[:]...)
}

func KeyRange() []byte {
	return []byte{'r', 'a', 'n', 'g', 'e'}
}

func WriteFlashLoanTx(db ethdb.KeyValueWriter, flashLoan *etypes.FlashLoan) {
	if err := db.Put(KeyFlashLoan(flashLoan.TXHash), flashLoan.Bytes()); err != nil {
		panic(fmt.Sprintf("error writing flash loan transaction to database: %s", err.Error()))
	}
}

func ReadFlashLoanTx(db ethdb.KeyValueStore, txHash common.Hash) *etypes.FlashLoan {
	data, _ := db.Get(KeyFlashLoan(txHash))
	if len(data) == 0 {
		return nil
	}

	flashLoan, err := etypes.DecodeFlashLoan(data)
	if err != nil {
		panic(fmt.Sprintf("flash loan transaction deserialization failed: %s", err.Error()))
	}

	return flashLoan
}

func ReadAllFlashLoan(db ethdb.KeyValueStore) []*etypes.FlashLoan {
	rst := []*etypes.FlashLoan{}
	it := db.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		if len(it.Key()) != flashLoanPrefixLength {
			continue
		}

		flashLoan, err := etypes.DecodeFlashLoan(it.Value())
		if err != nil {
			panic(fmt.Sprintf("flash loan transaction deserialization failed: %s", err.Error()))
		}
		rst = append(rst, flashLoan)
	}

	return rst
}

func DeleteFlashLoan(db ethdb.KeyValueStore, txHash common.Hash) {
	err := db.Delete(KeyFlashLoan(txHash))
	if err != nil {
		panic(fmt.Sprintf("failed to delete flash loan: %s", err.Error()))
	}
}

func WriteRange(db ethdb.KeyValueStore, start, limit uint64) {
	if err := db.Put(KeyRange(), append(EncodeUint64(start), EncodeUint64(limit)...)); err != nil {
		panic(fmt.Sprintf("an error occurred while writing the flash loan detection range to the database: %s", err.Error()))
	}
}

func ReadRange(db ethdb.KeyValueStore) (start, limit *uint64) {
	data, _ := db.Get(KeyRange())
	if len(data) != 8+8 {
		return
	}

	start = DecodeUint64(data[:8])
	limit = DecodeUint64(data[8:])
	return
}
