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
		values, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		words := strings.Split(values, " ")

		hermesKey := enigmas.Add(words[0], words[1])

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
