package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var Version = "dev"

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		log.Println("" +
			"This script accepts file paths in non-flag args which contain commands that will be executed (w/ bash -c '...'), one at a time, each after a prompt.",
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	log.Println("step-script @", Version)

	paths := flag.Args()

	if len(paths) == 0 {
		log.Fatalln("At least one path is required.")
	}

	commandsExecuted := 0
	for _, path := range paths {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalln(err)
		}
		scanner := bufio.NewScanner(bytes.NewReader(content))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
				continue
			}

			fmt.Print("\n$ ", line, "    # Execute? Y/n")
			if !yes(true) {
				log.Println("Skipping step...")
				continue
			}

			log.Println()

			command := exec.Command("bash", "-c", line)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err = command.Run()
			if err != nil {
				log.Fatal(err)
			}

			commandsExecuted++
		}
	}

	log.Printf("Finished running %d commands.", commandsExecuted)
}

func yes(byDefault bool) bool {
	response := strings.ToLower(strings.TrimSpace(prompt()))
	return (byDefault && response == "") || response == "y" || response == "yes"
}

func prompt() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}