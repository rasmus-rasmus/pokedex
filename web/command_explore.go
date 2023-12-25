package web

import (
	"encoding/json"
	"errors"
	"fmt"
)

func GetExploreAreaCallbackFct(conf *Config, client PokeAPIClient) func(cl_args []string) error {
	return func(cl_args []string) error {
		if len(cl_args) != 1 {
			return errors.New(fmt.Sprintf("Expected 1 command line argument. Got %v", len(cl_args)))
		}
		return exploreArea(cl_args[0], conf.Url, client)
	}
}

func exploreArea(location_name string, url string, client PokeAPIClient) error {
	fmt.Printf("Exploring %v...\n", location_name)
	fullUrl := makeFullUrl(url, location_name)
	body, err := client.fetchResource(fullUrl)
	if err != nil {
		return err
	}

	s := Location{}
	json.Unmarshal(body, &s)
	for _, enc := range s.PokemonEncounters {
		fmt.Printf("- %v\n", enc.Pokemon.Name)
	}
	return nil

}
