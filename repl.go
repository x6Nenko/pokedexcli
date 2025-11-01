package main

import (
    "bufio"
    "fmt"
    "strings"
		"os"
		"github.com/x6Nenko/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, client *pokeapi.Client) error
}

type config struct {
	next 		 *string
	previous *string
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
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
	}
}

func startRepl() {
	// Init client once that stays live for the entire session
	apiClient := pokeapi.NewClient("https://pokeapi.co/api/v2/location-area")
	cfg := &config{
		next:     nil,
		previous: nil,
	}

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
			err := cmd.callback(cfg, apiClient)  // Call the function stored in callback
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

func commandExit(cfg *config, client *pokeapi.Client) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp(cfg *config, client *pokeapi.Client) error {
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

func commandMap(cfg *config, client *pokeapi.Client) error {
    var urlToFetch string
    
    if cfg.next == nil {
        urlToFetch = "https://pokeapi.co/api/v2/location-area"
    } else {
        urlToFetch = *cfg.next
    }
    
    data, err := client.FetchItems(urlToFetch)

		if err != nil {
			fmt.Println(err)
			return err
		}

		cfg.next = data.Next
		cfg.previous = data.Previous

		for _, value := range data.Results {
			fmt.Println(value.Name)
		}

    return nil
}

func commandMapb(cfg *config, client *pokeapi.Client) error {
    var urlToFetch string
    
    if cfg.previous == nil {
				fmt.Println("you're on the first page")
        return nil
    } else {
        urlToFetch = *cfg.previous
    }
    
    data, err := client.FetchItems(urlToFetch)

		if err != nil {
			fmt.Println(err)
			return err
		}

		cfg.next = data.Next
		cfg.previous = data.Previous

		for _, value := range data.Results {
			fmt.Println(value.Name)
		}

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