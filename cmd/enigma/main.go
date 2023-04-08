package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/pipethedev/enigma"
)

func main() {
	create := flag.Bool("create", false, "Create a new connection key for your hermes connection")

	get := flag.Bool("get", false, "Get a connection key for your hermes connection")

	flag.Parse()

	enigmas := enigma.Enigmas{}

	scanner := bufio.NewScanner(os.Stdin)

	// Check if the user has sudo privileges
	cmd := exec.Command("sudo", "-n", "true")

	err := cmd.Run()
	if err != nil {
		println("This command requires sudo privileges.")
		return
	}

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
