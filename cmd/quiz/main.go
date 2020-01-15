package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func run() error {

	var csv string
	var limit int

	flag.StringVar(&csv, "csv", "problems.csv",
		"A CSV in the format of 'question,answer'")
	flag.IntVar(&limit, "limit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	return nil
}
