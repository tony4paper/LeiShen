package edb

import (
	"fmt"
	"leishen/pkg/etypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	platformNamePrefix    = byte(0x02)
	platformNamePrefixLen = 21
)

func KeyPlatformName(addr common.Address) []byte {
	return append([]byte{platformNamePrefix}, addr[:]...)
}

func WritePlatformName(db ethdb.KeyValueStore, platformName *etypes.PlatformName) {
	if err := db.Put(KeyPlatformName(platformName.Address), platformName.Bytes()); err != nil {
		panic(fmt.Errorf("an error occurred while writing the application name to the database: %v", err))
	}
}

func ReadPlatformName(db ethdb.KeyValueStore, addr common.Address) *etypes.PlatformName {
	data, err := db.Get(KeyPlatformName(addr))
	if err != nil {
		return nil
	}

	platformName, err := etypes.DecodePlatformName(data)
	if err != nil {
		panic(fmt.Errorf("an error occurred while deserializing the app name: %v", err))
	}

	return platformName
}

func DeleteAllPlatformName(db ethdb.KeyValueStore) {
	it := db.NewIterator(nil, nil)
	defer it.Release()

	for it.Next() {
		key := it.Key()
		if len(key) != platformNamePrefixLen {
			continue
		}

		if key[0] != platformNamePrefix {
			panic("the prefix of key mismatch with platform name")
		}

		db.Delete(key)
	}
}
