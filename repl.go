package main

import (
	"fmt"
	"strings"
	"os"
	"github.com/paul39-33/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name		string
	description	string
	callback	func(*config) error
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	//loop through all commands
	for _, innerMap := range commands {
		fmt.Printf("%s:%s\n", innerMap.name, innerMap.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	var mapRes pokeapi.LocationResponse
	var url string

	//check if its the first map call
	if cfg.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *cfg.Next
	}
	mapRes, err := pokeapi.FetchLocationAreas(url)
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

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
	} else {
		url := *cfg.Previous

		//Similar to commandMap
		mapRes, err := pokeapi.FetchLocationAreas(url)
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