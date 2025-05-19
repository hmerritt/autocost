package command

import (
	"fmt"
	"os"

	"github.com/hmerritt/go-car-monthly-cost-calculator/log"
	"github.com/hmerritt/go-car-monthly-cost-calculator/version"

	"github.com/mitchellh/cli"
)

func Run() {
	// Initiate new CLI app
	app := cli.NewCLI(version.AppName, version.GetVersion().VersionNumber())
	app.Args = os.Args[1:]

	// Setup logger
	logger := log.NewLogger()
	// go logger.FileStart(log.LOG_FILE)

	// Feed active commands to CLI app
	app.Commands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &RunCommand{
				BaseCommand: GetBaseCommand(logger),
			}, nil
		},
	}

	// Run app
	exitStatus, err := app.Run()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprint(err))
	}

	// logger.FileClose()

	// Exit without an error if no arguments were passed
	if len(app.Args) == 0 {
		os.Exit(0)
	}

	os.Exit(exitStatus)
}
