package predicate

import (
	"leishen/pkg/etypes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func GetEtherRecords(itx []*etypes.InternalTransaction) []*EtherRecord {
	ether_result := []*EtherRecord{}

	for _, record := range itx {
		ether_result = append(ether_result, &EtherRecord{record})
	}

	return ether_result
}

type EtherRecord struct {
	*etypes.InternalTransaction
}

func (r *EtherRecord) Token() string {
	panic("should not be used")
}

func (r *EtherRecord) From() string {
	panic("should not be used")
}

func (r *EtherRecord) To() string {
	panic("should not be used")
}

func (r *EtherRecord) Amount() *big.Int {
	return r.Value
}

func (r *EtherRecord) TokenAddress() common.Address {
	panic("should not be used")
}

func (r *EtherRecord) FromAddress() common.Address {
	panic("should not be used")
}

func (r *EtherRecord) ToAddress() common.Address {
	panic("should not be used")
}
