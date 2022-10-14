package setplatformname

import "github.com/urfave/cli/v2"

var (
	SetPlatformNameCommand = &cli.Command{
		Action: set,
		Name:   "set-platform-name",
		Usage:  "set as many addresses as possible with their platform names based on the platform names collected",
		Flags: []cli.Flag{
			contractDBFlag,
			platformNameDBFlag,
			csvFileFlag,
		},
	}
)

func set(c *cli.Context) error {
	contractDBPath := c.String(contractDBFlag.Name)
	contractDB, err := openDB(contractDBPath, 256, 64, true)
	if err != nil {
		return err
	}

	fltxDBPath := c.String(platformNameDBFlag.Name)
	fltxDB, err := openDB(fltxDBPath, 256, 64, false)
	if err != nil {
		return err
	}

	csvFilePath := c.String(csvFileFlag.Name)
	return setPlatformName(contractDB, fltxDB, csvFilePath)
}
