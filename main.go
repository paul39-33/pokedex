package main
import (
	"fmt"
	"bufio"
	"os"
)

func main(){

	commands := map[string]cliCommand{}
	
		commands["exit"] = cliCommand{
			name:		"exit",
			description:"Exit the Pokedex",
			callback:	commandExit,
		}
		commands["help"] = cliCommand{
			name:		"help",
			description:"Displays a help message",
			callback:	func() error {return commandHelp(commands)},
		}
	

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan(){
			text := scanner.Text()
			cleanText := cleanInput(text)
			userInput := cleanText[0]
			if command, exists := commands[userInput]; exists {
				err := command.callback()

				if err != nil {
					fmt.Println("Error: ", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error: ", err)
		}
	}
}


