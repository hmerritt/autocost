package command

import (
	"fmt"
	"strings"

	"github.com/hmerritt/autocost/ui"
	"github.com/hmerritt/autocost/utils"
)

type RunCommand struct {
	*BaseCommand
}

func (c *RunCommand) Synopsis() string {
	return "Run calculation"
}

func (c *RunCommand) Help() string {
	helpText := fmt.Sprintf(`
Usage: autocost run [options] TASK
  
Run calculation.

Example:
  $ autocost run
`)

	return strings.TrimSpace(helpText)
}

func (c *RunCommand) Flags() *FlagMap {
	return GetFlagMap(FlagNamesGlobal)
}

func (c *RunCommand) Run(args []string) int {
	// Gather auto information from user

	c.Log.Info("Information")
	// name := ui.AskString("Name: ")
	price := ui.AskFloat("Price")
	taxPerYear := ui.AskFloat("Tax per year")
	insurancePerYear := ui.AskFloat("Insurance per year")
	maintenancePerYear := ui.AskFloat("(Estimated) Maintenance per year")
	// milesDrivenPerYear := ui.AskFloat("Miles driven per year")
	// mpgAvg := ui.AskFloat("Average MPG")
	// fuelPricePerLiter := ui.AskFloat("Fuel price per liter")
	c.Log.Info("")

	// Calculate

	ownershipCostPerYearFor10Years := make([]float64, 0, 10)
	ownershipCostPerMonthFor10Years := make([]float64, 0, 10)
	ownershipCostPerDayFor10Years := make([]float64, 0, 10)

	for i := range 10 {
		ownershipLengthInYears := float64(i + 1)

		ac := &AutoCost{
			price:              price,
			taxPerYear:         taxPerYear,
			insurancePerYear:   insurancePerYear,
			maintenancePerYear: maintenancePerYear,
		}

		avgCostInYears := ac.Calc(ownershipLengthInYears)
		ownershipCostPerYearFor10Years = append(ownershipCostPerYearFor10Years, avgCostInYears)
		ownershipCostPerMonthFor10Years = append(ownershipCostPerMonthFor10Years, utils.FloatRound(avgCostInYears/12, 2))
		ownershipCostPerDayFor10Years = append(ownershipCostPerDayFor10Years, utils.FloatRound(avgCostInYears/365, 2))
	}

	fmt.Println("Cost per year for each year of ownership (1-10 years): ", ownershipCostPerYearFor10Years)
	fmt.Println("Cost per month for each year of ownership (1-10 years): ", ownershipCostPerMonthFor10Years)
	fmt.Println("Cost per day for each year of ownership (1-10 years): ", ownershipCostPerDayFor10Years)

	return 0
}

type AutoCost struct {
	price              float64
	taxPerYear         float64
	insurancePerYear   float64
	maintenancePerYear float64
	milesDrivenPerYear float64
	mpgAvg             float64
	fuelPricePerLiter  float64
}

func (a *AutoCost) Calc(ownershipLengthInYears float64) float64 {
	// Add up all the costs
	totalCost := a.price
	totalCost += a.taxPerYear * ownershipLengthInYears
	totalCost += a.insurancePerYear * ownershipLengthInYears
	totalCost += a.maintenancePerYear * ownershipLengthInYears

	// @TODO: Calculate fuel cost
	// fuelPricePerGallon := a.fuelPricePerLiter * 4.54609

	finalCostPerYear := totalCost / ownershipLengthInYears
	finalCostPerYear = utils.FloatRound(finalCostPerYear, 2)

	return finalCostPerYear
}
