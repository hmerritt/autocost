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

		fmt.Printf("Invalid input '%s'. Please enter 'y'/'yes' or 'n'/'no'.\n", normalizedInput)
	}
}
