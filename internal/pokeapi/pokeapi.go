package pokeapi

import (
	"encoding/json"
	"net/http"
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

func FetchLocationAreas(url string) (LocationResponse, error) {
	var locationRes LocationResponse
	res, err := http.Get(url)
	if err != nil {
		return locationRes, err
	}
	defer res.Body.Close()

	
	decoder := json.NewDecoder(res.Body)
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