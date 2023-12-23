package web

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Config struct {
	PrevResource int
	NextResource int
	Url          string
}

func GetNextMapCallbackFct(conf *Config, client PokeAPIClient) func(cl_args []string) error {
	return func(cl_args []string) error {
		if len(cl_args) > 0 {
			return errors.New("Too many command line arguments. Expected 0.")
		}
		return getNextLocations(conf, client)
	}
}

func GetPrevMapCallbackFct(conf *Config, client PokeAPIClient) func(cl_args []string) error {
	return func(cl_args []string) error {
		if len(cl_args) > 0 {
			return errors.New("Too many command line arguments. Expected 0.")
		}
		return getPrevLocations(conf, client)
	}
}

func GetExploreAreaFct(conf *Config, client PokeAPIClient) func(cl_args []string) error {
	return func(cl_args []string) error {
		if len(cl_args) != 1 {
			return errors.New(fmt.Sprintf("Expected 1 command line argument. Got %v", len(cl_args)))
		}
		return exploreArea(cl_args[0], conf.Url, client)
	}
}

func makeFullUrl(endPointUrl, resource string) string {
	if endPointUrl[0] != '/' {
		endPointUrl = "/" + endPointUrl
	}
	if endPointUrl[len(endPointUrl)-1] != '/' {
		endPointUrl = endPointUrl + "/"
	}
	return "https://pokeapi.co" + endPointUrl + resource + "/"
}

func exploreArea(location_name string, url string, client PokeAPIClient) error {
	fullUrl := makeFullUrl(url, location_name)
	body, err := client.fetchResource(fullUrl)
	if err != nil {
		return err
	}

	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Printf("Exploring %v\n", s.Name)
	for _, enc := range s.PokemonEncounters {
		fmt.Printf("- %v\n", enc.Pokemon.Name)
	}
	return nil

}

func getNextLocations(conf *Config, client PokeAPIClient) error {
	defer func(i, j *int) {
		*i += 20
		*j += 20
	}(&conf.PrevResource, &conf.NextResource)
	for i := 0; i < 20; i++ {
		nextResourceToFetch := conf.NextResource + i
		fmt.Printf("%v: ", nextResourceToFetch)
		err := getNextLocation(conf.Url, nextResourceToFetch, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPrevLocations(conf *Config, client PokeAPIClient) error {
	if conf.PrevResource < 1 {
		return errors.New("No more maps to show")
	}
	defer func(i, j *int) {
		*i -= 20
		*j -= 20
	}(&conf.PrevResource, &conf.NextResource)
	for i := 0; i < 20; i++ {
		nextResourceToFetch := conf.PrevResource + i
		fmt.Printf("%v: ", nextResourceToFetch)
		err := getPrevLocation(conf.Url, nextResourceToFetch, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func getNextLocation(url string, resource int, client PokeAPIClient) error {
	fullUrl := makeFullUrl(url, fmt.Sprint(resource))
	body, err := client.fetchResource(fullUrl)
	if err != nil {
		return err
	}

	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Name)
	return err
}

func getPrevLocation(url string, resource int, client PokeAPIClient) error {
	fullUrl := makeFullUrl(url, fmt.Sprint(resource))
	body, err := client.fetchResource(fullUrl)
	if err != nil {
		return err
	}

	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Name)
	return nil
}
