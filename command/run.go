package command

import (
	"fmt"
	"strings"
)

type RunCommand struct {
	*BaseCommand
}

func (c *RunCommand) Synopsis() string {
	return "Run calculation"
}

func (c *RunCommand) Help() string {
	helpText := fmt.Sprintf(`
Usage: go-car-monthly-cost-calculator run [options] TASK
  
Run calculation.

Example:
  $ go-car-monthly-cost-calculator run
`)

	return strings.TrimSpace(helpText)
}

func (c *RunCommand) Flags() *FlagMap {
	return GetFlagMap(FlagNamesGlobal)
}

func (c *RunCommand) Run(args []string) int {
	c.Log.Info("Command: run")

	return 0
}
