package web

import (
	"fmt"
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
