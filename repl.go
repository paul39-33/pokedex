package main

import (
	"fmt"
	"strings"
	"os"
	"github.com/paul39-33/pokedex/internal/pokeapi"
	"github.com/paul39-33/pokedex/internal/pokecache"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*config, *pokecache.Cache, []string, *pokeapi.Pokedex) error
}

type config struct {
	Previous	*string
	Next		*string
}


func cleanInput(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	//fmt.Println("Fields: ", fields)
	return fields
}

func commandExit(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, c *pokecache.Cache, args []string, commands map[string]cliCommand, p *pokeapi.Pokedex) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	//loop through all commands
	for _, innerMap := range commands {
		fmt.Printf("%s:%s\n", innerMap.name, innerMap.description)
	}
	return nil
}

func commandMap(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	var mapRes pokeapi.LocationResponse
	var url string

	//check if its the first map call
	if cfg.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *cfg.Next
	}
	mapRes, err := pokeapi.FetchLocationAreas(url, c)
	if err != nil {
		return fmt.Errorf("Error : %v", err)
	}
	//update config's Previous and Next url
	cfg.Next = mapRes.Next
	cfg.Previous = mapRes.Previous
	locations := pokeapi.GetLocations(mapRes)
	
	//prints the locations
	for _, location := range locations {
		fmt.Println(location)
	}
	return nil
}

func commandMapb(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
	} else {
		url := *cfg.Previous

		//Similar to commandMap
		mapRes, err := pokeapi.FetchLocationAreas(url, c)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}
		cfg.Next = mapRes.Next
		cfg.Previous = mapRes.Previous
		locations := pokeapi.GetLocations(mapRes)
		for _, location := range locations {
			fmt.Println(location)
		}
		
	}
	return nil
}

func commandExplore(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	if len(args) < 1 {
		return fmt.Errorf("Explore command requires an area name")
	}
	fmt.Printf("Exploring %v...\n", args[0])
	
	areaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)
	areaPokemon, err := pokeapi.FetchAreaPokemon(url, c)
	if err != nil {
		return fmt.Errorf("Error when trying to explore! Try checking the area name!")
	}
	pokeapi.GetPokemonList(areaPokemon)
	return nil
}

func commandCatch(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	if len(args) < 1 {
		return fmt.Errorf("Catch command requires a pokemon name")
	}

	if len(p.PokedexList) > 0 {
		_, exists := p.PokedexList[args[0]]
		if exists{
			return fmt.Errorf("%s already caught!", args[0])
		}
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0])
	pokemonInfo, err := pokeapi.GetPokemonInfo(url, c)
	if err != nil {
		return fmt.Errorf("Error trying to get pokemon information: %v", err)
	}
	err = pokeapi.CatchPokemon(pokemonInfo, p)
	if err != nil {
		return fmt.Errorf("Error catching pokemon: %v", err)
	}
	return nil
}

func commandInspect(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	pokemon, exists := p.PokedexList[args[0]]
	if exists {
		fmt.Printf("\nName: %v\n", pokemon.Name)
			fmt.Printf("Base Exp: %v\n", pokemon.BaseExp)
			fmt.Printf("Height: %v\n", pokemon.Height)
			fmt.Printf("Weight: %v\n", pokemon.Weight)
			fmt.Printf("Abilities: \n")
			for _, abilities := range pokemon.Abilities{
				fmt.Printf("\t-%v\n", abilities.Ability.Name)
			}
			fmt.Printf("Stats: \n")
			for _, stats := range pokemon.Stats{
				fmt.Printf("\t-%v: %v\n", stats.Stat.Name, stats.BaseStat)
			}
			fmt.Printf("Types: \n")
			for _, types := range pokemon.Types{
				fmt.Printf("\t-%v\n", types.Type.Name)
			}
			return nil
		}
	fmt.Println("You have not caught that pokemon")
	return nil
}


func commandPokedex(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {
	if len(p.PokedexList) > 0 {
		fmt.Println("Your Pokedex:")
		for _, pokemon := range p.PokedexList {
			fmt.Printf("\t- %v\n", pokemon.Name)
		}
		return nil
	}
	return fmt.Errorf("Pokedex is empty!\n")
}