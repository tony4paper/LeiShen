package edb

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"golang.org/x/crypto/sha3"

	"leishen/pkg/etypes"
)

const (
	keyPrefixLenthItx               = 40
	keyLenthItx                     = 48
	keyPrefixLengthCompletedItxFile = 35
)

var (
	completedItxFilePrefix = [...]byte{'i', 't', 'x'}
	completeItxFileFlag    = [...]byte{0x01}
)

func KeyItxSize(number uint64, txHash common.Hash) []byte {
	return KeyItxPrifix(number, txHash)
}

func KeyItx(number uint64, txHash common.Hash, index uint64) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, keyLenthItx))
	buf.Write(EncodeUint64(number))
	buf.Write(txHash[:])
	buf.Write(EncodeUint64(index))
	return buf.Bytes()
}

func KeyItxPrifix(number uint64, txHash common.Hash) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, keyPrefixLenthItx))
	buf.Write(EncodeUint64(number))
	buf.Write(txHash[:])
	return buf.Bytes()
}

func KeyCompletedItxFile(fileName string) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, keyPrefixLengthCompletedItxFile))

	buf.Write(completedItxFilePrefix[:])

	hasher := sha3.NewLegacyKeccak256()
	hasher.Reset()
	hasher.Write([]byte(fileName))
	hash := hasher.Sum(nil)
	buf.Write(hash)

	return buf.Bytes()
}

func WriteItxSize(db ethdb.KeyValueWriter, number uint64, txHash common.Hash, cnt uint64) {
	if err := db.Put(KeyItxSize(number, txHash), EncodeUint64(cnt)); err != nil {
		errMsg := fmt.Sprintf("An error occurred while writing the number of internal transactions to the database: %s", err.Error())
		panic(errMsg)
	}
}

func ReadItxSize(db ethdb.KeyValueReader, number uint64, txHash common.Hash) *uint64 {
	data, err := db.Get(KeyItxSize(number, txHash))
	if err != nil {
		return nil
	}
	return DecodeUint64(data)
}

func WriteItx(db ethdb.KeyValueWriter, tx *etypes.InternalTransaction) {
	if err := db.Put(KeyItx(tx.BlockNumber, tx.TransactionHash, tx.Index), tx.Bytes()); err != nil {
		errMsg := fmt.Sprintf("An error occurred while writing internal transactions to the database: %s", err.Error())
		panic(errMsg)
	}
}

func ReadItx(db ethdb.KeyValueReader, number uint64, txHash common.Hash, index uint64) *etypes.InternalTransaction {
	data, err := db.Get(KeyItx(number, txHash, index))
	if err != nil {
		return nil
	}

	tx, err := etypes.DecodeInternalTransaction(data)
	if err != nil {
		panic(fmt.Sprintf("An error occurred while deserializing the internal transaction: %s", err.Error()))
	}

	return tx
}

func ReadAllItxs(blockDB, itxDB ethdb.KeyValueStore, txHash common.Hash) []*etypes.InternalTransaction {
	number, _ := ReadIndexBytxHash(blockDB, txHash)
	if number == nil {
		return nil
	}

	return ReadAllItxsWithNumber(itxDB, *number, txHash)
}

func ReadAllItxsWithNumber(db ethdb.KeyValueStore, number uint64, txHash common.Hash) []*etypes.InternalTransaction {
	txs := []*etypes.InternalTransaction{}
	prefix := KeyItxPrifix(number, txHash)

	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == keyLenthItx {
			data := it.Value()
			tx, err := etypes.DecodeInternalTransaction(data)
			if err != nil {
				errMsg := fmt.Sprintf("Error reading internal transactions in batch: %s", err.Error())
				panic(errMsg)
			}
			txs = append(txs, tx)
		}
	}

	return txs
}

func WriteCompletedItxFileName(db ethdb.KeyValueWriter, fileName string) {
	if err := db.Put(KeyCompletedItxFile(filepath.Base(fileName)), completeItxFileFlag[:]); err != nil {
		panic(fmt.Sprintf("error writing to filename of completed internal transaction file: %v", err))
	}
}

func CheckCompletedItxFile(db ethdb.KeyValueReader, fileName string) bool {
	data, err := db.Get(KeyCompletedItxFile(filepath.Base(fileName)))
	return err == nil && len(data) == 1 && bytes.Equal(data, completeItxFileFlag[:])
}
