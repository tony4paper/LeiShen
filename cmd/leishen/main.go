package main

import (
	"fmt"
	"leishen/cmd/leishen/check"
	"leishen/cmd/leishen/predicate"
	setplatformname "leishen/cmd/leishen/set-platform-name"
	"leishen/cmd/leishen/statistics"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Usage = "leishen detects and analyzes flash loans to determine whether it is an attack against flash loans"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		check.FlashLoanCheckCommand,
		setplatformname.SetPlatformNameCommand,
		predicate.PredicateCommand,
		statistics.StatisticsCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
