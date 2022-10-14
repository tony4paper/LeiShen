package predicate

import "github.com/ethereum/go-ethereum/common"

type ResultWriter interface {
	Write(txHash common.Hash, number uint64, loanTypes, predicate string) error
	WriteError() error
}
