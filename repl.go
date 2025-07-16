package main

import (
	"fmt"
	"strings"
	"os"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}



func cleanInput(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	//fmt.Println("Fields: ", fields)
	return fields
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, innerMap := range commands {
		fmt.Printf("%s:%s\n", innerMap.name, innerMap.description)
	}
	return nil
}