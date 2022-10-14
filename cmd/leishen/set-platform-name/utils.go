package setplatformname

import (
	"strings"

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

func trimSliceSpace(rawRow []string) []string {
	var row []string
	emptyCnt := 0
	for _, s := range rawRow {
		ss := strings.TrimSpace(s)
		if len(ss) == 0 {
			emptyCnt++
		}
		row = append(row, ss)
	}

	if emptyCnt == csv_column_count {
		return nil
	}

	return row
}