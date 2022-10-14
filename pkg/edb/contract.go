package edb

import (
	"bytes"
	"fmt"
	"leishen/pkg/etypes"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"golang.org/x/crypto/sha3"
)

const (
	keyLenthContract                     = 20
	keyLenthSubContract                  = 40
	keyPrefixLengthCompletedContractFile = 34
)

var (
	subContractFlag             = [...]byte{0x01}
	completedContractFilePrefix = [...]byte{'k', 't'}
	completeContractFileFlag    = [...]byte{0x01}
)

func KeyContract(addr *common.Address) []byte {
	return addr[:]
}

func KeySubContract(creator, contract common.Address) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, keyLenthSubContract))
	buf.Write(creator[:])
	buf.Write(contract[:])
	return buf.Bytes()
}

func KeyCompletedContractFile(fileName string) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, keyPrefixLengthCompletedContractFile))

	buf.Write(completedContractFilePrefix[:])

	hasher := sha3.NewLegacyKeccak256()
	hasher.Reset()
	hasher.Write([]byte(fileName))
	hash := hasher.Sum(nil)
	buf.Write(hash)

	return buf.Bytes()
}

func WriteContract(db ethdb.KeyValueWriter, contract *etypes.Contract) {
	if err := db.Put(KeyContract(&contract.Address), contract.Bytes()); err != nil {
		panic(fmt.Errorf("error writing internal transaction to database: %v", err))
	}
}

func ReadContract(db ethdb.KeyValueReader, addr *common.Address) *etypes.Contract {
	data, err := db.Get(KeyContract(addr))
	if err != nil {
		return nil
	}

	contract, err := etypes.DecodeContract(data)
	if err != nil {
		panic(fmt.Errorf("an error occurred while deserializing the contract: %v", err))
	}
	return contract
}

func WriteSubContract(db ethdb.KeyValueWriter, contract *etypes.Contract) {
	if err := db.Put(KeySubContract(contract.Creator, contract.Address), subContractFlag[:]); err != nil {
		panic(fmt.Sprintf("error writing to filename of completed internal transaction file: %v", err))
	}
}

func ReadAllSubContract(db ethdb.KeyValueStore, addr common.Address) []*etypes.Contract {
	addrs := ReadAllSubContractAddress(db, addr)
	contracts := []*etypes.Contract{}
	for _, addr := range addrs {
		contracts = append(contracts, ReadContract(db, addr))
	}
	return contracts
}

func ReadAllSubContractAddress(db ethdb.KeyValueStore, addr common.Address) []*common.Address {
	addrs := []*common.Address{}
	iter := GetSubContractIterator(db, addr)
	defer iter.Release()
	for {
		addr, ok := iter.Next()
		if !ok {
			break
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

func GetSubContractIterator(db ethdb.KeyValueStore, addr common.Address) *subContractIterator {
	return &subContractIterator{db.NewIterator(addr[:], nil)}
}

type subContractIterator struct {
	iter ethdb.Iterator
}

func (iter *subContractIterator) Next() (*common.Address, bool) {
	if !iter.iter.Next() {
		return nil, false
	}

	key := iter.iter.Key()
	value := iter.iter.Value()
	if len(key) == keyLenthContract {
		return iter.Next()
	}

	if len(key) != keyLenthSubContract || len(value) != 1 || !bytes.Equal(value, subContractFlag[:]) {
		return nil, false
	}

	addr := common.BytesToAddress(key[keyLenthContract:])
	return &addr, true
}

func (iter *subContractIterator) Release() {
	iter.iter.Release()
}

func WriteCompletedContractFileName(db ethdb.KeyValueWriter, fileName string) {
	if err := db.Put(KeyCompletedContractFile(filepath.Base(fileName)), completeContractFileFlag[:]); err != nil {
		panic(fmt.Sprintf("error writing filename of completed contract file: %v", err))
	}
}

func CheckCompletedContractFile(db ethdb.KeyValueReader, fileName string) bool {
	data, err := db.Get(KeyCompletedContractFile(filepath.Base(fileName)))
	return err == nil && len(data) == 1 && bytes.Equal(data, completeContractFileFlag[:])
}
