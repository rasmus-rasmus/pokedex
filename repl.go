package main

import (
	"fmt"
	"os"
	"bufio"
	"errors"
	"strings"
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
	return errors.New("Exiting pokedex...")
}

func getCLIMap() map[string]cliCommand {

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
	}
}

func cleanInput(input string) []string {
	output := strings.ToLower(input)
	words := strings.Fields(output)
	return words
}


func tartRepl() {
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
				fmt.Println(err)
				break
			}
		}
	}
}