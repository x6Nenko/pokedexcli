package main

import (
    "bufio"
    "fmt"
    "strings"
		"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()           // Wait for user input
		input := scanner.Text()  // Get the input as a string
		cleanInput := cleanInput(input)

		if len(cleanInput) == 0 {
			continue
		}

		command := cleanInput[0]

		cmd, exists := getCommands()[command]  // Try to get the command from map
		if exists {                            // Did we find it?
			err := cmd.callback()       		 		 // Call the function stored in callback
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func cleanInput(text string) []string {
	// Convert to lowercase
	lower := strings.ToLower(text)

	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(lower)

	// Split by whitespace
	words := strings.Fields(trimmed)

	return words
}