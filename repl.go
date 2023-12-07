package main

import (
	"fmt"
	"os"
	"bufio"
	"errors"
	"strings"
	"web"
)


type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Print("Welcome to the pokedex!\n\n")

	cliMap := getCLIMap()

	if len(cliMap) == 0 {
		fmt.Println("No commands available!")
		return nil
	} else if cliMap == nil {
		return errors.New("Can't get cli map")
	}

	fmt.Println("You have the following options:")
	for _, v := range cliMap {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	fmt.Print("\n")
	 return nil 
}

func commandExit() error {
	return errors.New("exit")
}

func getCLIMap() map[string]cliCommand {
	c := web.Config{-39, 1, "https://pokeapi.co/api/v2/location-area/"}
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name: 		 "map",
			description: "Get next 20 locations",
			callback:	 web.GetNextMapCallbackFct(&c),
		},
		"mapb": {
			name: 		 "mapb",
			description: "Get previous 20 locations",
			callback: 	 web.GetPrevMapCallbackFct(&c),
		},
	}
}

func cleanInput(input string) []string {
	output := strings.ToLower(input)
	words := strings.Fields(output)
	return words
}


func startRepl() {
	cliMap := getCLIMap()
 
	scanner := bufio.NewScanner(os.Stdin) 
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if len(command) == 0 {
			continue
		}
		cmd, ok := cliMap[command[0]]
		if !ok {
			fmt.Println("Invalid command - type 'help' for options")
		} else {
			err := cmd.callback()
			if err != nil {
				if err.Error() == "exit" {
					break
				}
				fmt.Println(err)
				continue
			}
		}
	}
}