package sitemap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/roloum/gophercises/link"
)

//GetPages creates a sitemap for the
func GetPages(domainURL string, depth int, logger *log.Logger) ([]string,
	error) {

	if err := isValidURL(domainURL); err != nil {
		return nil, err
	}

	links, err := retrievePageLinks(domainURL)
	if err != nil {
		return nil, err
	}
	logger.Printf("%v\n", links)

	return []string{}, nil
}

//retrievePageLinks Retrieves the list of links in the URL page
func retrievePageLinks(pageURL string) ([]link.Link, error) {
	response, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	links, err := link.Parse(response.Body)
	for _, link := range links {
		fmt.Println(link)
	}

	return links, err
}

// TODO: validate domainURL
//isValidURL validates the format of the URL
func isValidURL(domainURL string) error {
	return nil
}
