package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("discomarkov", "Generates Markov chains from Discord chat exports")
	// command to export users from input files and turn the users into a json file containing all their messages across all files
	export := parser.NewCommand("export", "Exports users from input files and turn the users into a json file containing all their messages across all files")
	eInput := export.String("i", "input", &argparse.Options{Required: false, Help: "Input path", Default: "./exports/input"})
	eOutput := export.String("o", "output", &argparse.Options{Required: false, Help: "Output path", Default: "./exports/users.json"})
	// command to generate markov chains from the users.json file
	generate := parser.NewCommand("generate", "Generates Markov chains from the users.json file")
	gInput := generate.String("i", "input", &argparse.Options{Required: false, Help: "Input path", Default: "./exports/users.json"})
	gOutput := generate.String("o", "output", &argparse.Options{Required: false, Help: "Output path", Default: "./exports/defs.cl"})
	gFilters := generate.StringList("f", "filter", &argparse.Options{Required: false, Help: "Filter", Default: []string{}})
	// parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}
	if export.Happened() {
		fmt.Printf("Exporting users from %s to %s...\n", *eInput, *eOutput)
	} else if generate.Happened() {
		fmt.Printf("Generating markov chains from %s to %s...\n", *gInput, *gOutput)
		fmt.Printf("Filters: %v\n", *gFilters)
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}
}
