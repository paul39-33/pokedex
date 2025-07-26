package main
import (
	"fmt"
	"bufio"
	"os"
	"github.com/paul39-33/pokedex/internal/pokecache"
	"github.com/paul39-33/pokedex/internal/pokeapi"
	"time"
)

func main(){

	commands := map[string]cliCommand{}
	cfg := new(config)
	cache := pokecache.NewCache(5 * time.Minute)
	args := []string{}
	userPokedex := pokeapi.Pokedex{
		PokedexList: make(map[string]pokeapi.PokemonInfo),
	}
	
		commands["exit"] = cliCommand{
			name:		"exit",
			description:"Exit the Pokedex",
			callback:	commandExit,
		}
		commands["help"] = cliCommand{
			name:		"help",
			description:"Displays a help message",
			callback:	func(cfg *config, c *pokecache.Cache, args []string, p *pokeapi.Pokedex) error {return commandHelp(cfg, c, args,commands, p)},
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
		commands["explore"] = cliCommand{
			name:		"explore",
			description:"Explore map area",
			callback:	commandExplore,
		}
		commands["catch"] = cliCommand{
			name:		"catch",
			description:"Attempt to catch a pokemon",
			callback:	commandCatch,
		}
		commands["inspect"] = cliCommand{
			name:		"inspect",
			description:"Inspect a caught pokemon",
			callback:	commandInspect,
		}
		commands["pokedex"] = cliCommand{
			name:		"pokedex",
			description:"See list of caught pokemon",
			callback:	commandPokedex,
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
			if len(cleanText) > 1 {
				args = cleanText[1:]
			}
			
			if command, exists := commands[userInput]; exists {
				//use command callback to run the command based on user input
				err := command.callback(cfg, cache, args, &userPokedex)

				if err != nil {
					fmt.Println("Callback error: ", err)
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


