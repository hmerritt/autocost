package command

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aquasecurity/table"
	"github.com/hmerritt/autocost/ui"
	"github.com/hmerritt/autocost/utils"
	"github.com/hmerritt/autocost/version"
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

	name := ui.AskString("Name")
	price := ui.AskFloat("Price")
	taxPerYear := ui.AskFloat("Tax per year")
	insurancePerYear := ui.AskFloat("Insurance per year")
	maintenancePerYear := ui.AskFloat("(Estimated) Maintenance per year")
	milesDrivenPerYear := ui.AskFloat("(Estimated) Miles driven per year")
	mpgAvg := 0.0
	fuelPricePerLiter := 0.0
	if milesDrivenPerYear > 0 {
		mpgAvg = ui.AskFloat("Average MPG")
		fuelPricePerLiter = ui.AskFloat("Fuel price per liter")
	}
	c.Log.Info("")

	// Calculate

	yearsToCalculate := 10
	ownershipCostPerYearForNYears := make([]float64, 0, yearsToCalculate)
	ownershipCostPerMonthForNYears := make([]float64, 0, yearsToCalculate)
	ownershipCostPerDayForNYears := make([]float64, 0, yearsToCalculate)

	for i := range yearsToCalculate {
		ownershipLengthInYears := float64(i + 1)

		ac := &AutoCost{
			price:              price,
			taxPerYear:         taxPerYear,
			insurancePerYear:   insurancePerYear,
			maintenancePerYear: maintenancePerYear,
			milesDrivenPerYear: milesDrivenPerYear,
			mpgAvg:             mpgAvg,
			fuelPricePerLiter:  fuelPricePerLiter,
		}

		avgCostInYears := ac.Calc(ownershipLengthInYears)
		ownershipCostPerYearForNYears = append(ownershipCostPerYearForNYears, avgCostInYears)
		ownershipCostPerMonthForNYears = append(ownershipCostPerMonthForNYears, utils.FloatRound(avgCostInYears/12, 2))
		ownershipCostPerDayForNYears = append(ownershipCostPerDayForNYears, utils.FloatRound(avgCostInYears/365, 2))
	}

	// Log file

	logFile, _ := os.OpenFile(fmt.Sprintf("%s.log", version.AppName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer logFile.Close()
	logFileStat, err := logFile.Stat()
	if err == nil && logFileStat.Size() > 1 {
		logFile.WriteString("\n---\n\n")
	}
	stdPlusLog := io.MultiWriter(os.Stdout, logFile)

	// Print results

	logFile.WriteString(fmt.Sprintf(`Name: %s
Price: £%.2f
Tax per year: £%.2f
Insurance per year: £%.2f
(Estimated) Maintenance per year: £%.2f
Date of calculation: %s

`, name, price, taxPerYear, insurancePerYear, maintenancePerYear, time.Now().Format("Jan 2, 2006 at 3:04pm (MST)")))

	t := table.New(stdPlusLog)
	t.SetPadding(3)
	t.SetDividers(table.UnicodeRoundedDividers)
	// t.SetDividers(table.MarkdownDividers)
	t.SetHeaders("YEAR", "£/year", "£/month", "£/day")
	t.SetAlignment(table.AlignLeft, table.AlignRight, table.AlignRight, table.AlignRight)

	for y := range yearsToCalculate {
		t.AddRow(
			fmt.Sprint(y+1),
			fmt.Sprint(ownershipCostPerYearForNYears[y]),
			fmt.Sprint(ownershipCostPerMonthForNYears[y]),
			fmt.Sprint(ownershipCostPerDayForNYears[y]),
		)
	}

	t.Render()

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

	// Calculate fuel cost
	fuelPricePerGallon := a.fuelPricePerLiter * 4.54609
	fuelCost := (a.milesDrivenPerYear / a.mpgAvg) * fuelPricePerGallon * ownershipLengthInYears
	totalCost += fuelCost

	finalCostPerYear := totalCost / ownershipLengthInYears
	finalCostPerYear = utils.FloatRound(finalCostPerYear, 2)

	return finalCostPerYear
}
