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
	contractFiledNumber = reflect.TypeOf(Contract{}).NumField()
)

func init() {
	gob.Register(new(Contract))
}

func DecodeContract(data []byte) (*Contract, error) {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	contract := new(Contract)
	err := dec.Decode(contract)
	return contract, err
}

type Contract struct {
	Address           common.Address
	BlockNumber       uint64
	Timestamp         time.Time
	TransactionHash   common.Hash
	Creator           common.Address
	CreatorIsContract bool
	Value             *big.Int
	CreationCode      string
	ContractCode      string
}

func (c *Contract) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	enc.Encode(c)
	return buf.Bytes()
}

func (c *Contract) CSV() (string, error) {
	rst := make([]string, 0, contractFiledNumber)

	rst = append(rst, strings.ToLower(c.Address.Hex()))

	rst = append(rst, fmt.Sprint(c.BlockNumber))

	rst = append(rst, fmt.Sprint(c.Timestamp.Unix()))

	rst = append(rst, strings.ToLower(c.TransactionHash.Hex()))

	rst = append(rst, strings.ToLower(c.Creator.Hex()))

	if c.CreatorIsContract {
		rst = append(rst, "1")
	} else {
		rst = append(rst, "0")
	}

	rst = append(rst, c.Value.String())

	rst = append(rst, c.CreationCode)
	rst = append(rst, c.ContractCode)

	return strings.Join(rst, ","), nil
}
