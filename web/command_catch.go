package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func GetCatchPokemonCallbackFct(pokeMap map[string]Pokemon, client PokeAPIClient) func(cl_args []string) error {
	return func(cl_args []string) error {
		if len(cl_args) < 1 {
			return errors.New(fmt.Sprintf("Expected at least 1 command line argument. Got %v", len(cl_args)))
		}
		pokeballStrength := 50 // Regular pokeball
		if len(cl_args) == 1 {
			pokeballStrength = 50 // default if no ball type passed
		} else if cl_args[1] == "great" {
			pokeballStrength = 100
		} else if cl_args[1] == "ultra" {
			pokeballStrength = 150
		}
		return catchPokemon(cl_args[0], pokeballStrength, pokeMap, client)
	}
}

func catchPokemon(pokeName string, pokeballStrength int, pokeMap map[string]Pokemon, client PokeAPIClient) error {
	fullUrl := makeFullUrl("/api/v2/pokemon/", pokeName)
	pokeChannel := make(chan struct {
		pokemon Pokemon
		err     error
	})
	go func(ch chan struct {
		pokemon Pokemon
		err     error
	}) {
		pokemon := Pokemon{}
		body, err := client.fetchResource(fullUrl)
		if err == nil {
			json.Unmarshal(body, &pokemon)
		}
		ch <- struct {
			pokemon Pokemon
			err     error
		}{pokemon, err}
	}(pokeChannel)

	fmt.Printf("Throwing a pokeball at %v", pokeName)
	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Print("\n")

	pokeData := <-pokeChannel
	if pokeData.err != nil {
		return pokeData.err
	}

	if throwPokeball(pokeData.pokemon, pokeMap, -1) {
		fmt.Printf("%v was caught!\n", pokeName)
	} else {
		fmt.Printf("%v escaped!\n", pokeName)
	}
	return nil
}

// Set pokeballStrength to negative value to use default value of 50
func throwPokeball(pokemon Pokemon, pokeMap map[string]Pokemon, pokeballStrength int) bool {
	if pokeballStrength < 0 {
		pokeballStrength = 50
	}
	pokemonStrength := rand.Intn(pokemon.BaseExperience)
	if pokemonStrength > pokeballStrength {
		return false
	}
	pokeMap[pokemon.Name] = pokemon
	return true
}
