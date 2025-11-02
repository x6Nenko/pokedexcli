package main

import (
	"bufio"
	"fmt"
	"strings"
	"os"
	"github.com/x6Nenko/pokedexcli/internal/pokeapi"
	"time"
	"math/rand"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, cfg *config, client *pokeapi.Client) error
}

type config struct {
	next 		 *string
	previous *string
	pokedex  map[string]pokeapi.PokemonDetails
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
		"explore": {
			name:        "explore <area_name>",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "View details about a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View the list of caught Pokemons",
			callback:    commandPokedex,
		},
	}
}

func startRepl(interval time.Duration) {
	// Init client once that stays live for the entire session
	apiClient := pokeapi.NewClient(interval, "https://pokeapi.co/api/v2/location-area")

	cfg := &config{
		next:     nil,
		previous: nil,
		pokedex: make(map[string]pokeapi.PokemonDetails),
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
		args := cleanInput[1:]

		cmd, exists := getCommands()[command]  // Try to get the command from map
		if exists {                            // Did we find it?
			err := cmd.callback(args, cfg, apiClient)  // Call the function stored in callback
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

func commandExit(args []string, cfg *config, client *pokeapi.Client) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


func commandHelp(args []string, cfg *config, client *pokeapi.Client) error {
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

func commandMap(args []string, cfg *config, client *pokeapi.Client) error {
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

func commandMapb(args []string, cfg *config, client *pokeapi.Client) error {
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

func commandExplore(args []string, cfg *config, client *pokeapi.Client) error {
	if len(args) == 0 {
		fmt.Println("Usage: explore <area_name>")
		return nil
	}

	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")

	var urlToFetch = "https://pokeapi.co/api/v2/location-area/" + args[0]
	
	data, err := client.FetchLocationDetail(urlToFetch)

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, value := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", value.Pokemon.Name)
	}

	return nil
}

func commandCatch(args []string, cfg *config, client *pokeapi.Client) error {
	if len(args) == 0 {
		fmt.Println("Usage: catch <pokemon_name>")
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	var urlToFetch = "https://pokeapi.co/api/v2/pokemon/" + args[0]
	
	data, err := client.FetchPokemon(urlToFetch)

	if err != nil {
		fmt.Println(err)
		return err
	}

	pokemonName := data.Name
	baseExperience := data.BaseExperience

	// Generate a random number between 0-99
	roll := rand.Intn(100)

	// Calculate success threshold based on difficulty (baseExperience)
	// Lower threshold = harder to succeed
	successThreshold := 100 - (baseExperience / 5)

	catchResult := roll <= successThreshold

	if catchResult == true {
		// Checking if a pokemon already caught
		_, exists := cfg.pokedex[pokemonName]
		if exists {
			fmt.Printf("%s was already in your pokedex, go catch someone else...\n", pokemonName)
		} else {
			cfg.pokedex[pokemonName] = pokeapi.PokemonDetails{
				Name:           pokemonName,
				BaseExperience: baseExperience,
				Height:         data.Height,
				Weight:         data.Weight,
				Stats:          data.Stats,
				Types:          data.Types,
			}
			fmt.Printf("%s was caught!\n", pokemonName)
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(args []string, cfg *config, client *pokeapi.Client) error {
	if len(args) == 0 {
		fmt.Println("Usage: inspect <pokemon_name>")
		return nil
	}

	pokemonName := args[0]
	value, exists := cfg.pokedex[pokemonName]
	
	if exists {
		// direct access
		fmt.Printf("Name: %s\n", value.Name)
		fmt.Printf("Height: %d\n", value.Height)
		fmt.Printf("Weight: %d\n", value.Weight)
		fmt.Printf("Base Experience: %d\n", value.BaseExperience)
		
		// Stats - slice of structs, need to access nested Stat.Name
		fmt.Println("Stats:")
		for _, statEntry := range value.Stats {
			fmt.Printf("  -%s: %d\n", statEntry.Stat.Name, statEntry.BaseStat)
		}
		
		// Types - slice of structs, need to access nested Type.Name
		fmt.Println("Types:")
		for _, typeEntry := range value.Types {
			fmt.Printf("  - %s\n", typeEntry.Type.Name)
		}
	} else {
		fmt.Printf("You have not caught the %s yet ...\n", pokemonName)
	}

	return nil
}

func commandPokedex(args []string, cfg *config, client *pokeapi.Client) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your pokedex is empty...")
	}

	for _, pokemon := range cfg.pokedex {
		fmt.Printf("  - %s\n", pokemon.Name)
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