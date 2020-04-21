package link

import (
	"io"

	"golang.org/x/net/html"
)

//Link struct that holds the information for each a element in the HTML page
type Link struct {
	Href, Text string
}

//Parse parses the content of r and returns an array of Link with all the
//HTML a tags found in the document
func Parse(r io.Reader) ([]Link, error) {

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := []Link{}

	queue := []*html.Node{doc}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.Type == html.ElementNode && node.Data == "a" {

			var href string

			for _, a := range node.Attr {
				if a.Key == "href" {
					href = a.Val
					break
				}
			}

			links = append(links, Link{href, node.FirstChild.Data})

			node = node.FirstChild

		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			queue = append(queue, child)
		}
	}

	return links, nil
}
