package predicate

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

type TxHashIterator interface {
	Next() bool
	TxHash() *common.Hash
	Release()
}

func newTxHashIteratorFromStringSlice(txs []string) TxHashIterator {
	it := new(stringSliceIterator)
	for _, tx := range txs {
		txHash := common.HexToHash(tx)
		it.txHashes = append(it.txHashes, &txHash)
	}

	return it
}

type stringSliceIterator struct {
	txHashes []*common.Hash
	index    int
}

func (it *stringSliceIterator) Next() bool {
	it.index++
	return it.index <= len(it.txHashes)
}

func (it *stringSliceIterator) TxHash() *common.Hash {
	if it.index > len(it.txHashes) {
		return nil
	}
	return it.txHashes[it.index-1]
}

func (it *stringSliceIterator) Release() {}

func newTxHashIteratorFromFltxDB(fltxDB ethdb.Database) TxHashIterator {
	it := new(fltxDBIterator)
	it.iter = fltxDB.NewIterator(nil, nil)
	return it
}

type fltxDBIterator struct {
	iter ethdb.Iterator
	hash *common.Hash
}

func (it *fltxDBIterator) Next() bool {
	for it.iter.Next() {
		key := it.iter.Key()
		if len(key) != 33 {
			continue
		}

		hash := common.BytesToHash(key)
		it.hash = &hash

		return true
	}

	return it.iter.Next()
}

func (it *fltxDBIterator) TxHash() *common.Hash {
	return it.hash
}

func (it *fltxDBIterator) Release() {
	it.iter.Release()
}
