package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/anmitsu/go-shlex"
)

var Version = "dev"

func main() {
	log.SetFlags(0)
	path := flag.String("file", "", "The file with simple commands, one per line.")
	flag.Usage = func() {
		log.Println("" +
			"This script accepts individual commands as non-flag arguments " +
			"or in a text file and executes each after a confirmation (<ENTER>).",
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	commands := flag.Args()

	if len(commands) == 0 && *path != "" {
		raw, err := ioutil.ReadFile(*path)
		if err != nil {
			log.Fatal(err)
		}
		commands = strings.Split(string(raw), "\n")
	}

	if len(commands) == 0 {
		flag.Usage()
		log.Fatal("Please provide at least one command to run.")
	}

	log.Println("step-script @", Version)

	for _, command := range commands {
		line := strings.TrimSpace(command)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		args, err := shlex.Split(line, true)
		if err != nil {
			log.Fatal("Could not parse command:", command)
		}

		log.Println("\nAbout to run command:", line)
		fmt.Print("\n<ENTER> to continue")
		bufio.NewScanner(os.Stdin).Scan()

		log.Println()

		command := exec.Command(args[0], args[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Finished running %d commands.", len(commands))
}
