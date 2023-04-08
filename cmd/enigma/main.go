package main

import (
	"flag"

	"github.com/pipethedev/enigma"
)

func main() {
	create := flag.Bool("create", false, "Create a new connection key for your hermes connection")

	flag.Parse()

	enigmas := enigma.Enigmas{}

	switch {
	case *create:
		enigmas.Add("davmuri1414@gmail.com", "1234")
	default:
		println("No command specified.")
	}
}
