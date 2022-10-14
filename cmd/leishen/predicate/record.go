package predicate

import (
	"leishen/pkg/etypes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

func PreProcessRecord(contractDB, platformNameDB ethdb.Database, fltx *etypes.FlashLoan, record []Record) []Record {
	var result []Record
	result = RemoveZeroRecord(record)
	result = UnifiedAddressName(contractDB, platformNameDB, fltx, result)
	result = RemoveInternalRecord(result)
	result = MergeChange(result)
	result = RemoveWETHAndWBTC(result)
	result = MergeBlackHoleTransferred(result)
	result = MergeTransitTransferred(result)
	return result
}

type Record interface {
	Token() string
	From() string
	To() string
	Amount() *big.Int
	TokenAddress() common.Address
	FromAddress() common.Address
	ToAddress() common.Address
}

func WrapWithTokenName(record Record, token string) Record {
	return &RecordTokenNameWrapper{
		Inner:     record,
		TokenName: token,
	}
}

type RecordTokenNameWrapper struct {
	Inner     Record
	TokenName string
}

func (r *RecordTokenNameWrapper) Token() string {
	return r.TokenName
}

func (r *RecordTokenNameWrapper) From() string {
	return r.Inner.From()
}

func (r *RecordTokenNameWrapper) To() string {
	return r.Inner.To()
}

func (r *RecordTokenNameWrapper) Amount() *big.Int {
	return r.Inner.Amount()
}

func (r *RecordTokenNameWrapper) TokenAddress() common.Address {
	return r.Inner.TokenAddress()
}

func (r *RecordTokenNameWrapper) FromAddress() common.Address {
	return r.Inner.FromAddress()
}

func (r *RecordTokenNameWrapper) ToAddress() common.Address {
	return r.Inner.ToAddress()
}

func WrapWithPlatformName(record Record, from, to string) Record {
	return &RecordPlatformNameWrapper{
		Inner:    record,
		FromName: from,
		ToName:   to,
	}
}

type RecordPlatformNameWrapper struct {
	Inner    Record
	FromName string
	ToName   string
}

func (r *RecordPlatformNameWrapper) Token() string {
	return r.Inner.Token()
}

func (r *RecordPlatformNameWrapper) From() string {
	return r.FromName
}

func (r *RecordPlatformNameWrapper) To() string {
	return r.ToName
}

func (r *RecordPlatformNameWrapper) Amount() *big.Int {
	return r.Inner.Amount()
}

func (r *RecordPlatformNameWrapper) TokenAddress() common.Address {
	return r.Inner.TokenAddress()
}

func (r *RecordPlatformNameWrapper) FromAddress() common.Address {
	return r.Inner.FromAddress()
}

func (r *RecordPlatformNameWrapper) ToAddress() common.Address {
	return r.Inner.ToAddress()
}

func WrapWithChange(A, B, C Record, value *big.Int) Record {
	return &RecordChangeWrapper{
		A:     A,
		B:     B,
		C:     C,
		Value: value,
	}
}

type RecordChangeWrapper struct {
	A     Record
	B     Record
	C     Record
	Value *big.Int
}

func (r *RecordChangeWrapper) Token() string {
	return r.A.Token()
}

func (r *RecordChangeWrapper) From() string {
	return r.A.From()
}

func (r *RecordChangeWrapper) To() string {
	return r.A.To()
}

func (r *RecordChangeWrapper) Amount() *big.Int {
	return r.Value
}

func (r *RecordChangeWrapper) TokenAddress() common.Address {
	return r.A.TokenAddress()
}

func (r *RecordChangeWrapper) FromAddress() common.Address {
	return r.A.FromAddress()
}

func (r *RecordChangeWrapper) ToAddress() common.Address {
	return r.A.ToAddress()
}

func WrapWithBlackHole(A, B Record, value *big.Int) Record {
	return &RecordBlackHoleWrapper{
		A:     A,
		B:     B,
		Value: value,
	}
}

type RecordBlackHoleWrapper struct {
	A     Record
	B     Record
	Value *big.Int
}

func (r *RecordBlackHoleWrapper) Token() string {
	return r.A.Token()
}

func (r *RecordBlackHoleWrapper) From() string {
	return r.A.From()
}

func (r *RecordBlackHoleWrapper) To() string {
	return r.A.To()
}

func (r *RecordBlackHoleWrapper) Amount() *big.Int {
	return r.Value
}

func (r *RecordBlackHoleWrapper) TokenAddress() common.Address {
	return r.A.TokenAddress()
}

func (r *RecordBlackHoleWrapper) FromAddress() common.Address {
	return r.A.FromAddress()
}

func (r *RecordBlackHoleWrapper) ToAddress() common.Address {
	return r.A.ToAddress()
}

func WrapWithTransit(A, B Record, value *big.Int) Record {
	return &RecordTransitWrapper{
		A:     A,
		B:     B,
		Value: value,
	}
}

type RecordTransitWrapper struct {
	A     Record
	B     Record
	Value *big.Int
}

func (r *RecordTransitWrapper) Token() string {
	return r.A.Token()
}

func (r *RecordTransitWrapper) From() string {
	return r.A.From()
}

func (r *RecordTransitWrapper) To() string {
	return r.B.To()
}

func (r *RecordTransitWrapper) Amount() *big.Int {
	return r.Value
}

func (r *RecordTransitWrapper) TokenAddress() common.Address {
	return r.A.TokenAddress()
}

func (r *RecordTransitWrapper) FromAddress() common.Address {
	return r.A.FromAddress()
}

func (r *RecordTransitWrapper) ToAddress() common.Address {
	return r.B.ToAddress()
}

func WrapWithTrade(record Record, from, to string) Record {
	return &RecordTradeWrapper{
		Inner:    record,
		FromName: from,
		ToName:   to,
	}
}

type RecordTradeWrapper struct {
	Inner            Record
	FromName, ToName string
}

func (r *RecordTradeWrapper) Token() string {
	return r.Inner.Token()
}

func (r *RecordTradeWrapper) From() string {
	return r.FromName
}

func (r *RecordTradeWrapper) To() string {
	return r.ToName
}

func (r *RecordTradeWrapper) Amount() *big.Int {
	return r.Inner.Amount()
}

func (r *RecordTradeWrapper) TokenAddress() common.Address {
	return r.Inner.TokenAddress()
}

func (r *RecordTradeWrapper) FromAddress() common.Address {
	return r.Inner.FromAddress()
}

func (r *RecordTradeWrapper) ToAddress() common.Address {
	return r.Inner.ToAddress()
}
