package web

import (
	"errors"
	"fmt"
)

func GetInspectCallbackFct(pokeMap map[string]Pokemon) func(cl_args []string) error {
	return func(cl_args []string) error {
		if len(cl_args) != 1 {
			return errors.New(fmt.Sprintf("Expected 1 command line argument. Got %d", len(cl_args)))
		}
		pokemon, ok := pokeMap[cl_args[0]]
		if !ok {
			return errors.New("You have not caught that pokemon")
		}
		formatPokemonData(pokemon)
		return nil
	}
}

func formatPokemonData(pokemon Pokemon) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}
}
