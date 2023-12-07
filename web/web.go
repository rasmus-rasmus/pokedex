package web

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"errors"
)

type Config struct {
	PrevResource int
	NextResource int
	Url string
}

func GetNextMapCallbackFct(conf *Config) func() error {
	return func() error {
		return getNextLocations(conf)
	}
}

func GetPrevMapCallbackFct(conf *Config) func() error {
	return func() error {
		return getPrevLocations(conf)
	}
}

func getNextLocations(conf *Config) error {
	defer func (i *int) {
		*i += 20
	}(&conf.PrevResource)
	defer func (i *int) {
		*i += 20
	}(&conf.NextResource)
	for i := 0; i < 20; i++ {
		nextResourceToFetch := conf.NextResource + i
		fmt.Printf("%v: ", nextResourceToFetch)
		err := getNextLocation(conf.Url, nextResourceToFetch)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPrevLocations(conf *Config) error {
	if conf.PrevResource < 1 {
		return errors.New("No more maps to show")
	}
	defer func (i *int) {
		*i -= 20
	}(&conf.PrevResource)
	defer func (i *int) {
		*i -= 20
	}(&conf.NextResource)
	for i := 0; i < 20; i++ {
		nextResourceToFetch := conf.PrevResource + i
		fmt.Printf("%v: ", nextResourceToFetch)
		err := getPrevLocation(conf.Url, nextResourceToFetch)
		if err != nil {
			return err
		}
	}
	return nil
}

func getNextLocation(url string, resource int) error {
	res, err := http.Get(url+fmt.Sprint(resource)+"/")
	if err != nil {
		return errors.New("Failed to fetch resource")
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return errors.New("Response malformed")
	}
	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Request failed with status code %v", res.StatusCode))
	}
	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Name)
	return nil
}

func getPrevLocation(url string, resource int) error {
	res, err := http.Get(url+fmt.Sprint(resource)+"/")
	if err != nil {
		return errors.New("Failed to fetch resource")
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return errors.New("Request failed")
	}
	if err != nil {
		return errors.New("Response malformed")
	}
	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Name)
	return nil
}