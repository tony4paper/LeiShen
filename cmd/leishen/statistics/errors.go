package statistics

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func errDatabaseOpen(path string, err error) error {
	return fmt.Errorf("failed to open database (`%s`): %s", path, err.Error())
}

func errTxNotFount(txHash common.Hash) error {
	return fmt.Errorf("no tx found: %s", txHash)
}
