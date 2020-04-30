package sitemap

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/roloum/gophercises/link"
)

//GetPages creates a sitemap for the
func GetPages(domainURL string, depth int, logger *log.Logger) ([]string,
	error) {

	logger.Printf("Checking if %s is a valid domain\n", domainURL)
	if err := isValidURL(domainURL); err != nil {
		return nil, err
	}

	type node struct {
		page  string
		depth int
	}

	logger.Println("Domain added to the BFS queue")

	//depth starts on 1 because we are going to include in the pages array
	//the links of the last URL that we are going to process
	queue := []node{node{domainURL, 1}}

	pages := []string{fmt.Sprintf("%s/", domainURL)}

	visited := make(map[string]string)
	visited[fmt.Sprintf("%s/", domainURL)] = domainURL

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		//if node is beyond depth, skip
		if depth > 0 && n.depth > depth {
			continue
		}

		links, err := retrievePageLinks(n.page)
		if err != nil {
			return nil, err
		}

		for _, link := range links {

			var url string

			if strings.HasPrefix(link.Href, "http") {
				//Skip link if it's from different domain
				if !strings.HasPrefix(link.Href, domainURL) {
					continue
				}
				url = link.Href
			} else {
				//Add domain to link
				url = fmt.Sprintf("%s%s", domainURL, link.Href)
			}

			if _, ok := visited[url]; !ok {
				pages = append(pages, url)
				queue = append(queue, node{url, n.depth + 1})
				visited[url] = url
			}
		}

	}

	return pages, nil
}

//retrievePageLinks Retrieves the list of links in the URL page
func retrievePageLinks(pageURL string) ([]link.Link, error) {

	response, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return link.Parse(response.Body)
}

// TODO: validate domainURL
//isValidURL validates the format of the URL
func isValidURL(domainURL string) error {
	return nil
}
