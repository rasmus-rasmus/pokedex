package web

import (
	"cache"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PokeAPIClient struct {
	ClientCache *cache.Cache
}

func (client PokeAPIClient) fetchResource(resourceUrl string) ([]byte, error) {
	var body []byte
	val, ok := client.ClientCache.Get(resourceUrl)
	if ok {
		body = val
	} else {
		res, err := http.Get(resourceUrl)
		if err != nil {
			return []byte{}, errors.New("Failed to fetch resource")
		}
		resBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return []byte{}, errors.New("Response malformed")
		}
		if res.StatusCode > 299 {
			return []byte{}, errors.New(fmt.Sprintf("Request failed with status code %v", res.StatusCode))
		}
		body = resBody
		client.ClientCache.Add(resourceUrl, body)
	}
	return body, nil
}
