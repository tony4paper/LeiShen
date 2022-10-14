package check

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func errDatabaseOpen(path string, err error) error {
	return fmt.Errorf("failed to open database (`%s`): %s", path, err.Error())
}

func errBlockNotFount(number uint64) error {
	return fmt.Errorf("no block with height %d found", number)
}

func errTxNotFount(txHash common.Hash) error {
	return fmt.Errorf("no transaction found: %s", txHash)
}

func errReceiptTxNotFount(number uint64) error {
	return fmt.Errorf("could not find receipt for block %d", number)
}

func errBorrowerNotFound(hash common.Hash, name string) error {
	return fmt.Errorf("the transaction was identified as a flash loan on %s, but the borrower could not be found. Transaction hash:: %s", name, hash)
}
