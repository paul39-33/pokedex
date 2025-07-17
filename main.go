package main
import (
	"fmt"
	"bufio"
	"os"
	
)

func main(){

	commands := map[string]cliCommand{}
	cfg := new(config)
	
		commands["exit"] = cliCommand{
			name:		"exit",
			description:"Exit the Pokedex",
			callback:	commandExit,
		}
		commands["help"] = cliCommand{
			name:		"help",
			description:"Displays a help message",
			callback:	func(cfg *config) error {return commandHelp(cfg, commands)},
		}
		commands["map"] = cliCommand{
			name:		"map",
			description:"Displays list of location areas",
			callback:	commandMap,
		}
		commands["mapb"] = cliCommand{
			name:		"mapb",
			description:"Go back to previous maps page",
			callback:	commandMapb,
		}
	
	//setup the new scanner input method
	scanner := bufio.NewScanner(os.Stdin)
	//loop forever
	for {
		fmt.Print("Pokedex > ")
		//wait for user input
		if scanner.Scan(){
			text := scanner.Text()
			cleanText := cleanInput(text)
			userInput := cleanText[0]
			if command, exists := commands[userInput]; exists {
				//use command callback to run the command based on user input
				err := command.callback(cfg)

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


