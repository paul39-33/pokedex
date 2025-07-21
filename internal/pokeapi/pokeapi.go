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

type AreaPokemon struct {
	PokemonEncounters []struct{
		Pokemon Pokemon `json:"pokemon"`
		// You might also need to include other fields from this inner struct
		// if you intend to use them, e.g., VersionDetails
		// VersionDetails []interface{} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name	string `json:"name"`
	URL		string `json:"url"`
}

func FetchLocationAreas(url string, c *pokecache.Cache) (LocationResponse, error) {
	var locationRes LocationResponse

	if cache, found := c.Get(url); found {
		reader := bytes.NewReader(cache)

		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&locationRes); err != nil {
			return locationRes, err
		}
		return locationRes, nil
	}
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

func FetchAreaPokemon(url string, c *pokecache.Cache) (AreaPokemon, error) {
	var areaPokemonRes AreaPokemon

	if cache, found := c.Get(url); found {
		reader := bytes.NewReader(cache)

		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&areaPokemonRes); err != nil {
			return areaPokemonRes, fmt.Errorf("Error decoding cache : %v", err)
		}
		return areaPokemonRes, nil
	}

	res, err := http.Get(url)
	if err != nil  {
		return areaPokemonRes, fmt.Errorf("Error in http Get : %v", err)
	}
	defer res.Body.Close()

	//change res to []byte
	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return areaPokemonRes, fmt.Errorf("Error converting to []byte: %v", err)
	}
	c.Add(url, resByte)

	//change resByte to a readable format for decoder
	reader := bytes.NewReader(resByte)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&areaPokemonRes); err != nil {
		return areaPokemonRes, fmt.Errorf("Error decoding: %v", err)
	}
	return areaPokemonRes, nil
}

func GetPokemonList(res AreaPokemon){
	fmt.Println("Found Pokemon:")
	for _, pokemonList := range res.PokemonEncounters{
		fmt.Printf("- %v\n", pokemonList.Pokemon.Name)
	}
}