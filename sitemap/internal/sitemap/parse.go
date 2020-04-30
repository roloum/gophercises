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

	logger.Printf("Domain %s added to the BFS queue", domainURL)

	//depth starts on 1 because we are going to include in the pages array
	//the links of the last URL that we are going to process
	queue := []node{node{domainURL, 1}}

	pages := []string{domainURL}

	visited := make(map[string]string)
	visited[domainURL] = domainURL

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]

		//Since we are running BFS, once we find a depth beyond the parameter
		//we can break the loop
		if depth > 0 && n.depth > depth {
			logger.Printf(
				"Found link with depth(%d) beyond the parameter (%d)\nExit!",
				n.depth, depth)
			break
		}

		logger.Printf("Retrieving links for %s\n", n.page)
		links, err := retrievePageLinks(n.page)
		if err != nil {
			return nil, err
		}

		logger.Printf("Found %d links\n", len(links))

		added := 0
		for _, link := range links {

			var url string

			if strings.HasPrefix(link.Href, "http") {
				//Skip link if it's from different domain
				if !strings.HasPrefix(link.Href, domainURL) {
					continue
				}
				url = link.Href
			} else {
				//Omit # and mailto:
				if strings.HasPrefix(link.Href, "#") ||
					strings.HasPrefix(link.Href, "mailto:") ||
					//Found a link in the form
					//jon@calhoun.io
					//without mailto tag or /
					!strings.HasPrefix(link.Href, "/") {
					continue
				}

				//Add domain to link
				url = fmt.Sprintf("%s%s", domainURL, link.Href)
			}
			url = strings.TrimSuffix(url, "/")

			if _, ok := visited[url]; !ok {
				pages = append(pages, url)
				queue = append(queue, node{url, n.depth + 1})
				visited[url] = url
				added++
			}
		}
		logger.Printf("Added %d pages to the queue\n", added)
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
