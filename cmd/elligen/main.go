package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GodsBoss/elligen"
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
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	generator, err := elligen.DefaultGenerator()
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
