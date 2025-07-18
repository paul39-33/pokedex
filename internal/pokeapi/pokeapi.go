package pokeapi

import (
	"encoding/json"
	"net/http"
	"github.com/paul39-33/pokedex/internal/pokecache"
	"bytes"
	"fmt"
	"io"
)

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FetchLocationAreas(url string, c *pokecache.Cache) (LocationResponse, error) {
	var locationRes LocationResponse

	if cache, found := c.Get(url); found {
		fmt.Println("Cache found")
		reader := bytes.NewReader(cache)

		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&locationRes); err != nil {
			return locationRes, err
		}
		return locationRes, nil
	}
	fmt.Println("No cache found")
	res, err := http.Get(url)
	if err != nil {
		return locationRes, fmt.Errorf("GET error : %v", err)
	}
	defer res.Body.Close()

	//change res to []byte
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return locationRes, fmt.Errorf("Error converting to []byte: %v", err)
	}
	c.Add(url, resByte)

	//change resByte to a readable format for decoder
	reader := bytes.NewReader(resByte)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&locationRes); err != nil {
		return locationRes, err
	}
	return locationRes, nil
}

func GetLocations(res LocationResponse) []string {
	var locations []string
	for _, result := range res.Results{
		locations = append(locations, result.Name)
	}

	return locations
}