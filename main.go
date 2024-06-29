package main

import (
	"bufio"
	"fmt"
	"github.com/etiedem/pokedexcli/internal/pokeapi"
	"github.com/etiedem/pokedexcli/internal/pokecache"
	"github.com/etiedem/pokedexcli/internal/pokeconfig"
	"math/rand"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeconfig.Config, *string, *pokecache.Cache) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of all the pokemon in an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all of your pokemon",
			callback:    commandPokedex,
		},
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

func commandExplore(c *pokeconfig.Config, param *string, cache *pokecache.Cache) error {
	resp := pokeapi.SearchLocation(*param, cache)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *pokeconfig.Config, param *string, cache *pokecache.Cache) error {
	resp := get_pokemon(*param, cache)
	r := rand.Intn(resp.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...\n", resp.Name)
	if r > 40 {
		fmt.Printf("%s escaped!\n", resp.Name)
		return nil
	}
	fmt.Printf("%s was caught!\n", resp.Name)
	c.Pokedex[resp.Name] = resp

	return nil
}

func commandInspect(c *pokeconfig.Config, param *string, cache *pokecache.Cache) error {
	pokemon, exist := c.Pokedex[*param]

	if exist == false {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *pokeconfig.Config, _ *string, _ *pokecache.Cache) error {
	fmt.Println("Your Pokedex:")
	for name := range c.Pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandMap(c *pokeconfig.Config, param *string, cache *pokecache.Cache) error {
	resp := get_locations(c, "next", cache)
	for _, location := range resp.Locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(c *pokeconfig.Config, param *string, cache *pokecache.Cache) error {
	resp := get_locations(c, "previous", cache)
	for _, location := range resp.Locations {
		fmt.Println(location.Name)
	}
	return nil
}
func commandHelp(_ *pokeconfig.Config, _ *string, _ *pokecache.Cache) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit(_ *pokeconfig.Config, _ *string, _ *pokecache.Cache) error {
	os.Exit(0)
	return nil
}

func cleanInput(data string) []string {
	output := strings.ToLower(data)
	words := strings.Fields(output)
	return words
}

func main() {
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)

	config := &pokeconfig.Config{Pokedex: make(map[string]pokeconfig.Pokemon)}
	cache := pokecache.NewCache(time.Minute)

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		command, exists := commands[commandName]
		var param *string
		if len(words) > 1 {
			param = &words[1]
		}
		if exists {
			command.callback(config, param, cache)
		} else {
			fmt.Println("Unknown command")
		}
	}

}
