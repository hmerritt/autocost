//go:build windows

package main

import (
	"github.com/hmerritt/go-car-monthly-cost-calculator/command"
	"github.com/hmerritt/go-car-monthly-cost-calculator/version"
)

func main() {
	version.PrintTitle()
	command.Run()
}
