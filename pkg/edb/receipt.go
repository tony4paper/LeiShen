package edb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	keyReceipts             = []byte("receipts")
	keyLatestReceiptsNumber = keyLatestBlockNumber
)

func KeyReceipts(number uint64) []byte {
	return append(keyReceipts, EncodeUint64(number)...)
}

func KeyLatestReceiptsNumber() []byte {
	return keyLatestReceiptsNumber
}

func ReadRawReceipts(db ethdb.KeyValueStore, number uint64) types.Receipts {
	data, _ := db.Get(KeyReceipts(number))
	if len(data) == 0 {
		return nil
	}

	storageReceipts := []*types.ReceiptForStorage{}
	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
		panic(fmt.Sprintf("Invalid RLP receipt data. number: %d, err: %s", number, err.Error()))
	}
	receipts := make(types.Receipts, len(storageReceipts))
	for i, storageReceipt := range storageReceipts {
		receipts[i] = (*types.Receipt)(storageReceipt)
	}

	return receipts
}

func ReadReceiptsByBlock(receiptDB ethdb.KeyValueStore, block *types.Block) types.Receipts {
	receipts := ReadRawReceipts(receiptDB, block.NumberU64())
	if receipts == nil {
		return nil
	}

	if err := receipts.DeriveFields(params.MainnetChainConfig, block.Hash(), block.NumberU64(), block.Transactions()); err != nil {
		panic(fmt.Sprintf("An error occurred while filling the receipt with block data. number: %d, err: %s", block.NumberU64(), err.Error()))
	}
	return receipts
}

func ReadBlockAndReceipts(blockDB, receiptDB ethdb.KeyValueStore, number uint64) (*types.Block, types.Receipts) {
	block := ReadBlockByNumber(blockDB, number)
	if block == nil {
		return nil, nil
	}

	receipts := ReadRawReceipts(receiptDB, number)
	if receipts == nil {
		return block, nil
	}

	if err := receipts.DeriveFields(params.MainnetChainConfig, block.Hash(), number, block.Transactions()); err != nil {
		panic(fmt.Sprintf("An error occurred while filling the receipt with block data. number: %d, err: %s", number, err.Error()))
	}
	return block, receipts
}

func ReadReceipts(blockDB, receiptDB ethdb.KeyValueStore, number uint64) types.Receipts {
	_, receipts := ReadBlockAndReceipts(blockDB, receiptDB, number)
	return receipts
}

func ReadReceiptsByTxHash(blockDB, receiptDB ethdb.KeyValueStore, hash common.Hash) (types.Receipts, uint64) {
	number, index := ReadIndexBytxHash(blockDB, hash)
	if number == nil || index == nil {
		return nil, 0
	}

	return ReadReceipts(blockDB, receiptDB, *number), *index
}

func WriteReceipts(db ethdb.KeyValueWriter, number uint64, receipts types.Receipts) {
	storageReceipts := make([]*types.ReceiptForStorage, len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*types.ReceiptForStorage)(receipt)
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		panic(fmt.Sprintf("Receipt RLP code failed, err: %s", err.Error()))
	}

	if err := db.Put(KeyReceipts(number), bytes); err != nil {
		panic(fmt.Sprintf("An error occurred while writing the receipt to the receipt database, err: %s", err.Error()))
	}
}

func ReadLatestReceiptsNumber(db ethdb.KeyValueStore) *uint64 {
	data, _ := db.Get(keyLatestReceiptsNumber)
	if len(data) != 8 {
		return nil
	}

	return DecodeUint64(data)
}

func WriteLatestReceiptsNumber(db ethdb.KeyValueStore, number uint64) {
	if err := db.Put(keyLatestReceiptsNumber, EncodeUint64(number)); err != nil {
		panic(fmt.Sprintf("An error occurred while writing the latest receipt height to the receipt database. number:%d, err: %s", number, err.Error()))
	}
}
