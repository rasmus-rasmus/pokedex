package main

import (
	"bufio"
	"cache"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
	"web"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cl_args []string) error
}

func commandHelp(cl_args []string) error {
	fmt.Print("Welcome to the pokedex!\n\n")

	dummyCache, dummyChan := cache.NewCache(time.Duration(1 * time.Second))
	cliMap := getCLIMap(dummyCache)
	dummyChan <- struct{}{}

	if len(cliMap) == 0 {
		fmt.Println("No commands available!")
		return nil
	} else if cliMap == nil {
		return errors.New("Can't get cli map")
	}

	keys := make([]string, 0, len(cliMap))
	for k := range cliMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	fmt.Println("You have the following options:")
	for _, key := range keys {
		fmt.Printf("%v: %v\n", cliMap[key].name, cliMap[key].description)
	}
	fmt.Print("\n")
	return nil
}

func commandExit(cl_args []string) error {
	return errors.New("exit")
}

func getPokedexCallbackFct(pokeMap map[string]web.Pokemon) func(cl_args []string) error {
	return func(cl_args []string) error {
		allPokemon := make([]string, 0, len(pokeMap))
		for key := range pokeMap {
			allPokemon = append(allPokemon, key)
		}
		slices.Sort(allPokemon)
		fmt.Println("Your pokedex: ")
		for i, name := range allPokemon {
			fmt.Printf("  %d: %s\n", i+1, name)
		}
		return nil
	}
}

func getCLIMap(cache *cache.Cache) map[string]cliCommand {
	conf := web.Config{
		PrevResource: -39,
		NextResource: 1,
		Url:          "/api/v2/location-area/",
	}
	client := web.PokeAPIClient{ClientCache: cache}
	pokeMap := map[string]web.Pokemon{}
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
			description: "See next 20 locations",
			callback:    web.GetNextMapCallbackFct(&conf, client),
		},
		"mapb": {
			name:        "mapb",
			description: "See previous 20 locations",
			callback:    web.GetPrevMapCallbackFct(&conf, client),
		},
		"explore": {
			name:        "explore - expects one argument: location name",
			description: "Explore area",
			callback:    web.GetExploreAreaCallbackFct(&conf, client),
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon - expects one or two arguments: pokemon name and (optionally) pokeball to use; poke (default if none passed), great or ultra",
			callback:    web.GetCatchPokemonCallbackFct(pokeMap, client),
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a pokemon you have caught - expects one argument: pokemon name",
			callback:    web.GetInspectCallbackFct(pokeMap),
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all pokemon you have caught",
			callback:    getPokedexCallbackFct(pokeMap),
		},
	}
}

func cleanInput(input string) []string {
	output := strings.ToLower(input)
	words := strings.Fields(output)
	return words
}

func startRepl() {
	interval := time.Duration(20 * time.Second)
	cache, ch := cache.NewCache(interval)
	cliMap := getCLIMap(cache)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if len(command) == 0 {
			continue
		}
		cmd, ok := cliMap[command[0]]
		if !ok {
			fmt.Println("Invalid command - type 'help' for options")
		} else {
			err := cmd.callback(command[1:])
			if err != nil {
				if err.Error() == "exit" {
					ch <- struct{}{} // Telling reapLoop to return
					ch <- struct{}{} // Waiting for reapLoop to return
					break
				}
				fmt.Println(err)
				continue
			}
		}
	}
}
