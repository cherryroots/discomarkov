package main

import (
	"fmt"
	"log"
	"os"

	"discomarkov/internal"

	"github.com/akamensky/argparse"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("discomarkov", "Generates Markov chains from Discord chat exports")
	// command to export users from input files and turn the users into a json file containing all their messages across all files
	export := parser.NewCommand("export", "Exports users from input files and turn the users into a json file containing all their messages across all files")
	eInput := export.String("i", "input", &argparse.Options{Required: false, Help: "Input path", Default: "./exports/input"})
	eOutput := export.String("o", "output", &argparse.Options{Required: false, Help: "Output path", Default: "./exports/"})
	// command to generate markov chains from the users.json file
	generate := parser.NewCommand("generate", "Generates Markov chains from the users.json file")
	gInput := generate.String("i", "input", &argparse.Options{Required: false, Help: "Input file", Default: "./exports/users.json"})
	gOutput := generate.String("o", "output", &argparse.Options{Required: false, Help: "Output path", Default: "./exports"})
	// uid:1234567890, u:cherry, rid:1234567890, r:admin
	gFilters := generate.StringList("f", "filter", &argparse.Options{Required: false, Help: "Filter", Default: []string{}})
	gOutType := generate.String("t", "type", &argparse.Options{Required: false, Help: "Output type", Default: "json"})
	// parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		log.Println(parser.Usage(err))
		os.Exit(1)
	}
	if export.Happened() {
		log.Printf("Exporting users from %s to %s...\n", *eInput, *eOutput)
		exports, err := internal.ParseAllFiles(*eInput)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		users, err := internal.CollectUsers(exports)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		internal.WriteUsers(users, *eOutput+"users.json")
		log.Println("Exported users to", *eOutput+"users.json")
	} else if generate.Happened() {
		log.Printf("Generating markov chains from %s to %s...\n", *gInput, *gOutput)
		log.Printf("Filters: %v\n", *gFilters)
		users, err := internal.ReadUsers(*gInput)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		userMarkovs, err := internal.GenerateMarkovChains(users, *gFilters)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if *gOutType == "clips" {
			deffacts, err := internal.UserMarkovsToClips(userMarkovs)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			internal.WriteClipsDeffacts(deffacts, *gOutput)
		} else {
			internal.WriteUserMarkovs(userMarkovs, *gOutput+"."+*gOutType)
		}
		log.Printf("Generated markov chains to %s as %s\n", *gOutput, *gOutType)
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		log.Println(parser.Usage(err))
		os.Exit(1)
	}
}
