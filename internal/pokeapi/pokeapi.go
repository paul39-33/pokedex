package pokeapi

import (
	"encoding/json"
	"net/http"
	"github.com/paul39-33/pokedex/internal/pokecache"
	"bytes"
	"fmt"
	"io"
	//"math"
	"math/rand"
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

type PokemonInfo struct {
	Name		string `json:"name"`
	BaseExp		int `json:"base_experience"`
	Height		int `json:"height"`
	Weight		int `json:"weight"`
	Abilities	[]struct {
		Ability Ability `json:"ability"`
	} `json:"abilities"`
	Stats		[]struct{
		BaseStat	int `json:"base_stat"`
		Stat 		Stat `json:"stat"`
	} `json:"stats"`
	Types		[]struct{
		Type Type `json:"type"`
	} `json:"types"`
}

type Ability struct {
	Name	string `json:"name"`
	URL		string `json:"url"`
}

type Stat struct {
	Name	string `json:"name"`
	URL		string `json:"url"`
}

type Type struct {
	Name	string `json:"name"`
	URL		string `json:"url"`
}

type Pokedex struct {
	PokedexList map[string]PokemonInfo
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

func GetPokemonInfo(url string, c *pokecache.Cache) (PokemonInfo, error) {
	var pokemonInfo PokemonInfo

	if cache, found := c.Get(url); found {
		reader := bytes.NewReader(cache)

		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&pokemonInfo); err != nil {
			return pokemonInfo, fmt.Errorf("Error decoding cache: %v", err)
		}
		return pokemonInfo, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return pokemonInfo, fmt.Errorf("Pokemon Info not found: %v", err)
	}
	defer res.Body.Close()

	//check if res's StatusCode means pokemon not found
	if res.StatusCode != 200 {
		return pokemonInfo, fmt.Errorf("Pokemon not found!")
	}

	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return pokemonInfo, fmt.Errorf("Error converting to byte: %v", err)
	}
	c.Add(url, resByte)

	reader := bytes.NewReader(resByte)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&pokemonInfo); err != nil {
		return pokemonInfo, fmt.Errorf("Error decoding: %v", err)
	}
	return pokemonInfo, nil
}

func CatchPokemon(pokeInfo PokemonInfo, p *Pokedex) error {
	baseExp := pokeInfo.BaseExp / 100
	//base chance is 25% but gets smaller the higher the pokemon's base experience is
	chance := 4 + baseExp
	if rand.Intn(chance) == 3 {
		fmt.Printf("%s was caught!\n", pokeInfo.Name)
		p.PokedexList[pokeInfo.Name] = pokeInfo
		return nil
	}
	fmt.Printf("%s escaped!\n", pokeInfo.Name)
	return nil
}