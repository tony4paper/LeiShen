package etypes

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
)

func init() {
	gob.Register(new(FlashLoan))
	gob.Register(new(TransferredToken))
}

func DecodeFlashLoan(data []byte) (*FlashLoan, error) {
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	tx := new(FlashLoan)
	err := dec.Decode(tx)
	return tx, err
}

type FlashLoan struct {
	TXHash       common.Hash
	Types2Loaner map[string]common.Address
}

func (r *FlashLoan) Types() []string {
	types := make([]string, 0, len(r.Types2Loaner))
	for t := range r.Types2Loaner {
		types = append(types, t)
	}
	return types
}

func (r *FlashLoan) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	enc.Encode(r)
	return buf.Bytes()
}

type TransferredToken struct {
	Name      string
	From      common.Address
	FromName  string
	To        common.Address
	ToName    string
	Value     string
	TokenName string
	RawHtml   string
}
