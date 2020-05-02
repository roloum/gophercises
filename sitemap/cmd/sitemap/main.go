package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/roloum/gophercises/sitemap/internal/sitemap"
)

const appName = "sitemap"

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func main() {

	if err := run(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	var domainURL, logFileName string
	var depth int

	flag.StringVar(&domainURL, "domainURL", "https://www.calhoun.io",
		"URL used to build the site map")
	flag.IntVar(&depth, "depth", 0, "Maximum number of links to follow")
	flag.StringVar(&logFileName, "logFileName", "/dev/null", "Log file name")
	flag.Parse()

	logf, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	logger := log.New(logf, appName,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	if domainURL == "" {
		return errors.New("Domain can not be empty")
	} else if depth < 0 {
		return errors.New("Depth can not be negative")
	}

	logger.Println("Configuration loaded")
	logger.Println(domainURL, depth)

	pages, err := sitemap.GetPages(domainURL, depth, logger)
	if err != nil {
		return err
	}

	/*
		pagesXML, err := xml.Marshal(pages)
		for _, p := range pagesXML {
			fmt.Println(string(p))
		}
	*/
	type loc struct {
		Value string `xml:"loc"`
	}
	type urlset struct {
		Urls  []loc  `xml:"url"`
		Xmlns string `xml:"xmlns,attr"`
	}
	toXml := urlset{Xmlns: xmlns}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	fmt.Println(xmlns)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		return err
	}
	fmt.Println()

	return nil
}
