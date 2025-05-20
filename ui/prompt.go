package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ask user for y/n confirmation.
func AskConfirmation(prompt string, keepPromptingOnInvalidInput bool) bool {
	// Padding
	fmt.Println("")
	defer fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	// Loop until we get a valid response
	for {
		fmt.Printf("%s [y/n]: ", prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
			return false
		}

		normalizedInput := strings.ToLower(strings.TrimSpace(input))

		if normalizedInput == "y" || normalizedInput == "yes" {
			return true
		}
		if normalizedInput == "n" || normalizedInput == "no" {
			return false
		}

		if !keepPromptingOnInvalidInput {
			return false
		}

		fmt.Printf("Invalid input. Please enter 'y'/'yes' or 'n'/'no'.\n", normalizedInput)
	}
}

func AskFloat(prompt string) float64 {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s: ", prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
			return 0
		}

		normalizedInput := strings.TrimSpace(input)
		normalizedInput = strings.ReplaceAll(normalizedInput, ",", "")

		if normalizedInput == "" {
			continue
		}

		var value float64
		n, err := fmt.Sscanf(normalizedInput, "%f", &value)
		if n != 1 || err != nil {
			fmt.Printf("Invalid input. Please enter a valid number.\n", normalizedInput)
			continue
		}

		return value
	}
}

func AskString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s: ", prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
			return ""
		}

		normalizedInput := strings.TrimSpace(input)

		if normalizedInput == "" {
			continue
		}

		return normalizedInput
	}
}
