package edb

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	keyBlockHeader       = []byte("header")
	keyBlockBody         = []byte("body")
	keyLatestBlockNumber = []byte("latest")
)

func GetBlockNumberKey(blockHash common.Hash) []byte {
	return blockHash.Bytes()
}

func GetHeaderKey(number uint64) []byte {
	return append(keyBlockHeader, EncodeUint64(number)...)
}

func GetBodyKey(number uint64) []byte {
	return append(keyBlockBody, EncodeUint64(number)...)
}

func GetTxHashKey(hash common.Hash) []byte {
	return hash.Bytes()
}

func GetLatestBlockNumberKey() []byte {
	return keyLatestBlockNumber
}

func ReadBlockByBlockHash(db ethdb.KeyValueStore, hash common.Hash) *types.Block {
	number := ReadBlockNumber(db, hash)
	if number == nil {
		return nil
	}

	return ReadBlockByNumber(db, *number)
}

func ReadBlockNumber(db ethdb.KeyValueStore, hash common.Hash) *uint64 {
	data, _ := db.Get(GetBlockNumberKey(hash))
	if len(data) != 8 {
		return nil
	}

	return DecodeUint64(data)
}

func ReadBlockByNumber(db ethdb.KeyValueStore, number uint64) *types.Block {
	header := ReadHeader(db, number)
	if header == nil {
		return nil
	}

	body := ReadBody(db, number)
	if body == nil {
		return nil
	}

	return types.NewBlockWithHeader(header).WithBody(body.Transactions, body.Uncles)
}

func ReadHeader(db ethdb.KeyValueStore, number uint64) *types.Header {
	data, _ := db.Get(GetHeaderKey(number))
	if len(data) == 0 {
		return nil
	}

	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		panic(fmt.Sprint("Invalid chunk header RLP data read from database", "number:", number, "err:", err))
	}

	return header
}

func ReadBody(db ethdb.KeyValueStore, number uint64) *types.Body {
	data, _ := db.Get(GetBodyKey(number))
	if len(data) == 0 {
		return nil
	}

	body := new(types.Body)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		panic(fmt.Sprint("The chunk RLP data read from the database is invalid", "number:", number, "err:", err))
	}

	return body
}

func WriteBlock(db ethdb.KeyValueWriter, block *types.Block) {
	WriteHeader(db, block.Hash(), block.NumberU64(), block.Header())
	WriteBody(db, block.Hash(), block.NumberU64(), block.Body())
	WriteTxs(db, block)
}

func WriteHeader(db ethdb.KeyValueWriter, hash common.Hash, number uint64, header *types.Header) {

	WriteBlockNumber(db, hash, number)

	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		panic(fmt.Sprint("RLP encoding of block header data failed", "err:", err))
	}

	key := GetHeaderKey(number)
	if err := db.Put(key, data); err != nil {
		panic(fmt.Sprint("Failed to write block header data to the database", "err:", err))
	}
}

func WriteBody(db ethdb.KeyValueWriter, hash common.Hash, number uint64, body *types.Body) {

	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		panic(fmt.Sprint("RLP encoding of block data failed", "err:", err))
	}

	if err := db.Put(GetBodyKey(number), data); err != nil {
		panic(fmt.Sprint("Failed to write chunk data to database", "err:", err))
	}
}

func WriteBlockNumber(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	key := GetBlockNumberKey(hash)
	value := EncodeUint64(number)
	if err := db.Put(key, value); err != nil {
		panic(fmt.Sprint("Failed to write block height to database", "err:", err))
	}
}

func ReadBlockByTxHash(db ethdb.KeyValueStore, hash common.Hash) (*types.Block, *uint64) {
	number, index := ReadIndexBytxHash(db, hash)
	if number == nil || index == nil {
		return nil, nil
	}

	return ReadBlockByNumber(db, *number), index
}

func ReadIndexBytxHash(db ethdb.KeyValueStore, hash common.Hash) (*uint64, *uint64) {
	data, _ := db.Get(GetTxHashKey(hash))
	if len(data) != 16 {
		return nil, nil
	}

	number := DecodeUint64(data[:8])
	index := DecodeUint64(data[8:])
	return number, index
}

func ReadTxByTxHash(db ethdb.KeyValueStore, hash common.Hash) *types.Transaction {
	block, index := ReadBlockByTxHash(db, hash)
	if block == nil || index == nil {
		return nil
	}

	return block.Transactions()[*index]
}

func WriteTxs(db ethdb.KeyValueWriter, block *types.Block) {
	numberBytes := EncodeUint64(block.NumberU64())
	for i, tx := range block.Transactions() {
		WriteTx(db, tx.Hash(), numberBytes, EncodeUint64(uint64(i)))
	}
}

func WriteTx(db ethdb.KeyValueWriter, hash common.Hash, numberBytes, index []byte) {
	if err := db.Put(GetTxHashKey(hash), append(numberBytes, index...)); err != nil {
		panic(fmt.Sprint("Failed to write transaction to database", "err:", err))
	}
}

func ReadLatestBlockNumber(db ethdb.KeyValueStore) *uint64 {
	data, _ := db.Get(GetLatestBlockNumberKey())
	if len(data) != 8 {
		return nil
	}

	return DecodeUint64(data)
}

func WriteLatestBlockNumber(db ethdb.KeyValueStore, number uint64) {
	if err := db.Put(GetLatestBlockNumberKey(), EncodeUint64(number)); err != nil {
		panic(fmt.Sprint("Failed to write the latest block height to the database", "err:", err))
	}
}
