package setplatformname

import (
	"fmt"
)

func errDatabaseOpen(path string, err error) error {
	return fmt.Errorf("failed to open database (`%s`): %s", path, err.Error())
}

func errCSVOpen(err error) error {
	return fmt.Errorf("failed to open csv: %s", err.Error())
}

func errCSVRead(file_name string, err error) error {
	return fmt.Errorf("error reading csv file <%s>: %s", file_name, err.Error())
}

func errCSVItemNumber(expect, actual int) error {
	return fmt.Errorf("the amount of data is not as expected. Expected: %d. Actual: %d", expect, actual)
}
