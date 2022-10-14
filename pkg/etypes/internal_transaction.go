package etypes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

var (
	itxFiledNumber = reflect.TypeOf(InternalTransaction{}).NumField()
)

func init() {
	gob.Register(new(InternalTransaction))
}

func DecodeInternalTransaction(data []byte) (*InternalTransaction, error) {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	tx := new(InternalTransaction)
	err := dec.Decode(tx)
	return tx, err
}

type InternalTransaction struct {
	Index           uint64
	BlockNumber     uint64
	Timestamp       time.Time
	TransactionHash common.Hash
	CallIndex       string
	From            common.Address
	To              common.Address
	FromIsContract  bool
	ToIsContract    bool
	Value           *big.Int
	CallingFunction *string
	Error           *string
}

func (tx *InternalTransaction) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	enc.Encode(tx)
	return buf.Bytes()
}

func (tx *InternalTransaction) CSV() string {
	rst := make([]string, 0, itxFiledNumber)

	rst = append(rst, fmt.Sprint(tx.BlockNumber))

	rst = append(rst, fmt.Sprint(tx.Timestamp.Unix()))

	rst = append(rst, strings.ToLower(tx.TransactionHash.Hex()))

	rst = append(rst, tx.CallIndex)

	rst = append(rst, strings.ToLower(tx.From.Hex()))

	rst = append(rst, strings.ToLower(tx.To.Hex()))

	if tx.FromIsContract {
		rst = append(rst, "1")
	} else {
		rst = append(rst, "0")
	}

	if tx.ToIsContract {
		rst = append(rst, "1")
	} else {
		rst = append(rst, "0")
	}

	rst = append(rst, tx.Value.String())

	if tx.CallingFunction == nil {
		rst = append(rst, "0x")
	} else {
		rst = append(rst, "0x"+*tx.CallingFunction)
	}

	if tx.Error == nil {
		rst = append(rst, "None")
	} else {
		rst = append(rst, *tx.Error)
	}

	return strings.Join(rst, ",")
}
