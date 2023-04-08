package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/pipethedev/enigma"
)

func main() {
	create := flag.Bool("create", false, "Create a new connection key for your hermes connection")

	get := flag.Bool("get", false, "Get a connection key for your hermes connection")

	flag.Parse()

	enigmas := enigma.Enigmas{}

	scanner := bufio.NewScanner(os.Stdin)

	switch {
	case *create:

		fmt.Print("What is your email address ? ")

		scanner.Scan()
		email := scanner.Text()

		fmt.Print("Provide an app key for encryption ? ")

		scanner.Scan()
		key := scanner.Text()

		hermesKey := enigmas.Add(email, key)

		fmt.Printf("Hermes Key: %s", hermesKey)
		break
	case *get:
		fmt.Print("What is your email address ? ")

		scanner.Scan()
		email := scanner.Text()

		hermesKey := enigmas.Get(email)

		fmt.Printf("Hermes Key: %s", hermesKey)

	default:
		println("No command specified.")
	}
}
