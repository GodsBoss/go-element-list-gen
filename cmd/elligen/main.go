package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	elligen "github.com/GodsBoss/go-element-list-gen"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("elements name generator", flag.ContinueOnError)

	attempts := fs.Int("attempts", 0, "Number of failed attempts when generating names.")
	count := fs.Int("count", 10, "Number of names to be generated.")

	consonants := fs.String("consonants", "", "Comma-separated list of consonants. A default list is used if this is empty.")
	vowels := fs.String("vowels", "", "Comma-separated list of vowels. A default list is used if this is empty.")
	endings := fs.String("endings", "", "Comma-separated list of endings. A default list is used if this is empty.")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	options := make([]elligen.DefaultGeneratorOption, 0)

	if len(*consonants) > 0 {
		options = append(options, elligen.WithConsonants(strings.Split(*consonants, ",")))
	}
	if len(*vowels) > 0 {
		options = append(options, elligen.WithVowels(strings.Split(*vowels, ",")))
	}
	if len(*endings) > 0 {
		options = append(options, elligen.WithEndings(strings.Split(*endings, ",")))
	}

	generator, err := elligen.DefaultGenerator(options...)
	if err != nil {
		return err
	}

	list, err := elligen.Generate(generator, *count, *attempts)
	if len(list) > 0 {
		fmt.Println("Generated names:")
		for i := range list {
			fmt.Printf("- %s\n", list[i])
		}
	}
	return err
}
