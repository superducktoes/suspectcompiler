package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type packageInfo struct {
	Name        string
	Version     string
	Edition     string
	Description string
	Homepage    string
}

type tomlConfig struct {
	PI packageInfo `toml:"packageInfo"`
}

func parseTomlFile(conf *tomlConfig) {

	// read the toml file for config data
	content, err := ioutil.ReadFile("./config.toml")

	// check to see if there are any errors
	if err != nil {
		log.Fatal(err)
	}

	// convert to a string
	tomlText := string(content)

	if _, err := toml.Decode(tomlText, &conf); err != nil {
		log.Fatal(err)
	}

	fmt.Println("suspectcompiler, a markdown compiler by nick roy")
	fmt.Printf("Version: %s\n", conf.PI.Version)

}

func parseMarkdownFile(fileName string, tokens *[]string) {

	fmt.Printf("Attempting to open %s\n", fileName)

	// attempt to open the file from the command line
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Reading file...")
	}

	// create a new scanner to start reading line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var txtlines []string

	// store each line separate in the array
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	// read through our array and start parsing each line to create our tokens
	for _, line := range txtlines {

		var firstChar string
		var outputString string

		// first we check to see if we have an empty line or not
		if len(line) == 0 {
			outputString = "<br>"
		} else {
			firstChar = line[0:1]
			if firstChar == "#" {
				outputString = "<h1>" + line[2:] + "</h1>"
			} else {
				// anything else is going to be a normal paragraph
				outputString = "<p>" + line[0:] + "</p>"
			}
		}
		*tokens = append(*tokens, outputString)
	}

}

func main() {

	// get a list of command line arguments
	argsWithProg := os.Args
	var fileName string

	// make sure we just have the one argument
	if len(argsWithProg) != 2 {
		log.Fatal("Incorrect number of command line arguments")
	} else {
		fileName = argsWithProg[1]
	}

	// create a new object for reading our config file
	var conf tomlConfig

	// declare an empty slice for storing our tokens when parsing lines
	var tokens []string

	parseTomlFile(&conf)
	parseMarkdownFile(fileName, &tokens)
}
