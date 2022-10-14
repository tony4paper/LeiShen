package predicate

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	Erc20TopicFeature = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	Erc20TopicLength  = 1 + 2
)

func GetErc20Records(receipt *types.Receipt) []*Erc20Record {
	records := []*Erc20Record{}
	for _, log := range receipt.Logs {
		if len(log.Topics) == Erc20TopicLength && log.Topics[0] == Erc20TopicFeature {
			from := common.BytesToAddress(log.Topics[1].Bytes())
			to := common.BytesToAddress(log.Topics[2].Bytes())
			value := new(big.Int).SetBytes(log.Data)

			records = append(records, &Erc20Record{
				TokenAddr: log.Address,
				FromAddr:  from,
				ToAddr:    to,
				Value:     value,
			})
		}
	}

	return records
}

type Erc20Record struct {
	TokenAddr common.Address
	FromAddr  common.Address
	ToAddr    common.Address
	Value     *big.Int
}

func (erc20 *Erc20Record) Token() string {
	return erc20.TokenAddr.Hex()
}

func (erc20 *Erc20Record) From() string {
	return erc20.FromAddr.Hex()
}

func (erc20 *Erc20Record) To() string {
	return erc20.ToAddr.Hex()
}

func (erc20 *Erc20Record) Amount() *big.Int {
	return erc20.Value
}

func (erc20 *Erc20Record) TokenAddress() common.Address {
	return erc20.TokenAddr
}

func (erc20 *Erc20Record) FromAddress() common.Address {
	return erc20.FromAddr
}

func (erc20 *Erc20Record) ToAddress() common.Address {
	return erc20.ToAddr
}
