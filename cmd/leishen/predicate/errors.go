package predicate

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func errDatabaseOpen(path string, err error) error {
	return fmt.Errorf("failed to open database (`%s`): %s", path, err.Error())
}

func errTxNotFount(txHash common.Hash) error {
	return fmt.Errorf("no transaction found: %s", txHash)
}

func errReceiptTxNotFount(number uint64) error {
	return fmt.Errorf("could not find receipt for block %d", number)
}

func errTxNotFlashLoan(txHash common.Hash) error {
	return fmt.Errorf("%s is not a flash loan transaction", txHash)
}

func errTxToMessage(txhash common.Hash, err error) error {
	return fmt.Errorf("failed to convert transaction %s to message: %s", txhash, err.Error())
}
