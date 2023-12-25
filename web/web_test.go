package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestMakeFullUrl(t *testing.T) {
	cases := []struct {
		endPointUrl     string
		resource        string
		expectedFullUrl string
	}{
		{
			"/api/v2/location-area/",
			"1",
			"https://pokeapi.co/api/v2/location-area/1/",
		},
		{
			"api/v2/location-area/",
			"1",
			"https://pokeapi.co/api/v2/location-area/1/",
		},
		{
			"/api/v2/location-area",
			"1",
			"https://pokeapi.co/api/v2/location-area/1/",
		},
		{
			"api/v2/location-area",
			"1",
			"https://pokeapi.co/api/v2/location-area/1/",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			fullUrl := makeFullUrl(c.endPointUrl, c.resource)
			if fullUrl != c.expectedFullUrl {
				t.Errorf("Expected %v but got %v", c.expectedFullUrl, fullUrl)
			}
		})
	}
}

func TestCatchPokemon(t *testing.T) {
	cases := []struct {
		pokeName               string
		shouldBeAddedToPokedex bool
	}{
		{
			"pikachu",
			false,
		},
		{
			"pikachu",
			true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			fullUrl := makeFullUrl("/api/v2/pokemon/", c.pokeName)
			res, _ := http.Get(fullUrl)
			body, _ := io.ReadAll(res.Body)
			res.Body.Close()
			pokemon := Pokemon{}
			json.Unmarshal(body, &pokemon)
			var pokeballStrength int
			if c.shouldBeAddedToPokedex {
				pokeballStrength = pokemon.BaseExperience
			} else {
				pokeballStrength = 0
			}
			pokeMap := make(map[string]Pokemon)
			throwPokeball(pokemon, pokeMap, pokeballStrength)
			_, ok := pokeMap[c.pokeName]
			if ok != c.shouldBeAddedToPokedex {
				t.Errorf("%v was added: %v. Should have been added: %v", c.pokeName, ok, c.shouldBeAddedToPokedex)
			}
		})
	}
}
