package main

import (
    "bufio"
    "fmt"
    "strings"
		"os"
)

func cleanInput(text string) []string {
	// Convert to lowercase
	lower := strings.ToLower(text)

	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(lower)

	// Split by whitespace
	words := strings.Fields(trimmed)

	return words
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()           // Wait for user input
		input := scanner.Text()  // Get the input as a string
		cleanInput := cleanInput(input)
		fmt.Println("Your command was:", cleanInput[0])
	}
}