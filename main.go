package main

import (
	"fmt"
	"bufio"
	"os"
	"errors"
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


func main() {
	cliMap := getCLIMap()

	for {
		fmt.Print("Pokedex > ")
		var f *os.File 
		f = os.Stdin 
		defer f.Close() 
	 
		scanner := bufio.NewScanner(f) 
		scanner.Scan()
		command := scanner.Text()
		cmd, ok := cliMap[command]
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