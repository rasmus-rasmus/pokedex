package web

import (
	"cache"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	PrevResource int
	NextResource int
	Url          string
}

func GetNextMapCallbackFct(conf *Config, cache *cache.Cache) func() error {
	return func() error {
		return getNextLocations(conf, cache)
	}
}

func GetPrevMapCallbackFct(conf *Config, cache *cache.Cache) func() error {
	return func() error {
		return getPrevLocations(conf, cache)
	}
}

func getNextLocations(conf *Config, cache *cache.Cache) error {
	defer func(i *int) {
		*i += 20
	}(&conf.PrevResource)
	defer func(i *int) {
		*i += 20
	}(&conf.NextResource)
	for i := 0; i < 20; i++ {
		nextResourceToFetch := conf.NextResource + i
		fmt.Printf("%v: ", nextResourceToFetch)
		err := getNextLocation(conf.Url, nextResourceToFetch, cache)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPrevLocations(conf *Config, cache *cache.Cache) error {
	if conf.PrevResource < 1 {
		return errors.New("No more maps to show")
	}
	defer func(i *int) {
		*i -= 20
	}(&conf.PrevResource)
	defer func(i *int) {
		*i -= 20
	}(&conf.NextResource)
	for i := 0; i < 20; i++ {
		nextResourceToFetch := conf.PrevResource + i
		fmt.Printf("%v: ", nextResourceToFetch)
		err := getPrevLocation(conf.Url, nextResourceToFetch, cache)
		if err != nil {
			return err
		}
	}
	return nil
}

func getNextLocation(url string, resource int, cache *cache.Cache) error {
	var body []byte
	fullUrl := url + fmt.Sprint(resource) + "/"
	val, ok := cache.Get(fullUrl)
	if ok {
		body = val
	} else {
		res, err := http.Get(fullUrl)
		if err != nil {
			return errors.New("Failed to fetch resource")
		}
		resBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return errors.New("Response malformed")
		}
		if res.StatusCode > 299 {
			return errors.New(fmt.Sprintf("Request failed with status code %v", res.StatusCode))
		}
		body = resBody
		cache.Add(fullUrl, body)
	}
	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Name)
	return nil
}

func getPrevLocation(url string, resource int, cache *cache.Cache) error {
	var body []byte
	fullUrl := url + fmt.Sprint(resource) + "/"
	val, ok := cache.Get(fullUrl)
	if ok {
		body = val
	} else {
		res, err := http.Get(url + fmt.Sprint(resource) + "/")
		if err != nil {
			return errors.New("Failed to fetch resource")
		}
		resBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return errors.New("Request failed")
		}
		if err != nil {
			return errors.New("Response malformed")
		}
		body = resBody
		cache.Add(fullUrl, body)
	}
	s := Location{}
	json.Unmarshal(body, &s)
	fmt.Println(s.Name)
	return nil
}
