package predicate

import (
	"math/big"

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

// rf = (v - base) / base
func rateFluctuation(v, base *big.Float) *big.Float {
	c := big.NewFloat(0).Sub(v, base)
	return c.Quo(c, base)
}
