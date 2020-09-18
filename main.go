package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/anmitsu/go-shlex"
)

var Version = "dev"

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Println("" +
			"This script accepts individual commands as non-flag arguments " +
			"and executes each after a confirmation (<ENTER>).",
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		log.Fatal("Please provide at least one command to run.")
	}

	log.Println("step-script @ ", Version)

	for _, command := range flag.Args() {
		line := strings.TrimSpace(command)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		args, err := shlex.Split(line, true)
		if err != nil {
			log.Fatal("Could not parse command:", command)
		}

		fmt.Println("About to run command:")
		fmt.Println("\t" + line)
		fmt.Print("<ENTER> to continue")
		bufio.NewScanner(os.Stdin).Scan()

		command := exec.Command(args[0], args[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
