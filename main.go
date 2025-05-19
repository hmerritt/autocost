//go:build windows

package main

import (
	"github.com/hmerritt/autocost/command"
	"github.com/hmerritt/autocost/version"
)

func main() {
	version.PrintTitle()
	command.Run()
}
