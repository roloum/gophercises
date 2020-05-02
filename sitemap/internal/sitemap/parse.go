package sitemap

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/roloum/gophercises/link"
)

//GetPages creates a sitemap for the
func GetPages(domainURL string, depth int, logger *log.Logger) ([]string,
	error) {
	return bfs(domainURL, depth, logger)
}

//bfs performs the Breath First Search on the given URL and returns
//an array with all the links
func bfs(domainURL string, depth int, logger *log.Logger) ([]string,
	error) {

	type node struct {
		page  string
		depth int
	}

	logger.Printf("Domain %s added to the BFS queue", domainURL)

	//depth starts on 1 because we are going to include in the pages array
	//the links of the last URL that we are going to process
	queue := []node{node{domainURL, 1}}

	pages := []string{domainURL}

	visited := make(map[string]struct{})
	visited[domainURL] = struct{}{}

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
		urls, err := getPage(n.page, logger)
		if err != nil {
			return nil, err
		}

		logger.Printf("Found %d links\n", len(urls))

		added := 0
		for _, url := range urls {

			if _, ok := visited[url]; !ok {
				pages = append(pages, url)
				queue = append(queue, node{url, n.depth + 1})
				visited[url] = struct{}{}
				added++
			}
		}
		logger.Printf("Added %d pages to the queue\n", added)
	}

	return pages, nil
}

//getPage Retrieves the list of links in the URL page
func getPage(pageURL string, logger *log.Logger) ([]string, error) {

	logger.Printf("Requesting URL %s\n", pageURL)

	response, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	baseURL := &url.URL{
		Scheme: response.Request.URL.Scheme,
		Host:   response.Request.URL.Host,
	}

	return getLinks(response.Body, baseURL.String())
}

func getLinks(r io.Reader, baseURL string) ([]string, error) {
	links, err := link.Parse(r)
	if err != nil {
		return nil, err
	}

	urls := []string{}

	for _, link := range links {
		var url string

		switch {
		case strings.HasPrefix(link.Href, "/"):
			url = baseURL + link.Href
		case strings.HasPrefix(link.Href, "http"):
			url = link.Href
		}

		urls = append(urls, strings.TrimSuffix(url, "/"))
	}

	return filter(urls, withPrefix(baseURL)), nil
}

func filter(urls []string, filterFunc func(string) bool) []string {
	result := []string{}

	for _, url := range urls {
		if filterFunc(url) {
			result = append(result, url)
		}
	}

	return result
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

// TODO: validate domainURL
//isValidURL validates the format of the URL
func isValidURL(domainURL string) error {
	return nil
}
