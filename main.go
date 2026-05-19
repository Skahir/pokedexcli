package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	conf := &pokeConfig{
		Next: "",
		Prev: "",
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		dirty := scanner.Text()
		clean := cleanInput(dirty)

		commands := getCommands()

		val, ok := commands[clean[0]]
		if ok == true {
			err := val.callback(conf)

			if err != nil {
				fmt.Println("Something went wrong at the end")
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func commandExit(conf *pokeConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *pokeConfig) error {
	commands := getCommands()
	fmt.Printf(`Welcome to the Pokedex!
Usage:

`)

	for _, val := range commands {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	fmt.Println()

	return nil
}

func commandMap(conf *pokeConfig) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	if conf.Next != "" {
		url = conf.Next
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("something went wrong in the get for areas")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	m := MapArea{}

	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("something went wrong unmarshalling the mapArea")
	}

	conf.Next = m.Next
	conf.Prev = m.Previous

	for _, val := range m.Results {
		fmt.Println(val.Name)
	}

	return nil
}

func commandMapb(conf *pokeConfig) error {
	if conf.Prev == "" {
		fmt.Println("you're on the first page")
		conf.Next = "https://pokeapi.co/api/v2/location-area/"
	} else {
		resp, err := http.Get(conf.Prev)
		if err != nil {
			fmt.Println("mapb get pokedex data error")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		m := MapArea{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			fmt.Println("unmarshal in mapb error")
		}

		conf.Next = m.Next
		conf.Prev = m.Previous

		for _, val := range m.Results {
			fmt.Println(val.Name)
		}
	}

	return nil
}

func cleanInput(text string) []string {
	if len(text) <= 0 {
		return []string{}
	}

	clean := strings.TrimSpace(strings.Replace(text, "  ", " ", -1))
	clean = strings.ToLower(clean)
	split := strings.Split(clean, " ")

	return split
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Shows 20 names of location areas in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows the previous 20 names of location areas in the pokemon world",
			callback:    commandMapb,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(conf *pokeConfig) error
}

type pokeConfig struct {
	Next string
	Prev string
}

type MapArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
