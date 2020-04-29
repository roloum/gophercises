package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/roloum/gophercises/sitemap/internal/sitemap"
)

var appName = "sitemap"

func main() {

	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	var domainURL string
	var depth int

	logger := log.New(os.Stdout, appName,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	logger.Println("Parsing configuration")

	flag.StringVar(&domainURL, "domainURL", "https://www.calhoun.io",
		"URL used to build the site map")
	flag.IntVar(&depth, "depth", 0, "Maximum number of links to follow")
	flag.Parse()

	if domainURL == "" {
		return errors.New("Domain can not be empty")
	} else if depth < 0 {
		return errors.New("Depth can not be negative")
	}

	logger.Println("Configuration loaded")
	logger.Println(domainURL, depth)

	urls, err := sitemap.GetPages(domainURL, depth, logger)
	if err != nil {
		return err
	}
	logger.Printf("\n%s", strings.Join(urls, "\n"))

	return nil
}
