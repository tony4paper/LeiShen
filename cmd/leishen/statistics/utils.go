package statistics

import (
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
)

func openDB(path string, cache int, handles int, readonly bool) (ethdb.Database, error) {
	leveldb, err := leveldb.New(path, cache, handles, path, readonly)
	if err != nil {
		return nil, errDatabaseOpen(path, err)
	}

	db := rawdb.NewDatabase(leveldb)
	return db, nil
}
