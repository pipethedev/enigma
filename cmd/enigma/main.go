package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pipethedev/enigma"
)

func main() {
	create := flag.Bool("create", false, "Create a new connection key for your hermes connection")

	flag.Parse()

	enigmas := enigma.Enigmas{}

	switch {
	case *create:
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("What is your email address ? ")

		scanner.Scan()
		email := scanner.Text()

		fmt.Print("Provide an app key for encryption ? ")

		scanner.Scan()
		key := scanner.Text()

		hermesKey := enigmas.Add(email, key)

		fmt.Printf("Hermes Key: %s", hermesKey)
	default:
		println("No command specified.")
	}
}

func getInput(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty data is not allowed")
	}

	return text, nil

}
