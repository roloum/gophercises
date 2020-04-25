package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link struct that holds the information for each a element in the HTML page
type Link struct {
	Href, Text string
}

//Parse parses the content of r using breath first approach and returns an
// array of Link with all the HTML a tags found in the document
func Parse(r io.Reader) ([]Link, error) {

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return getLinks(doc), nil

}

//Searches through the entire tree of the html using breath first approach
//and returns an array with the "a href" elements found in the document
func getLinks(root *html.Node) []Link {
	links := []Link{}
	queue := []*html.Node{root}

	for len(queue) > 0 {

		//Dequeue
		node := queue[0]
		queue = queue[1:]

		if node.Type == html.ElementNode && node.Data == "a" {
			links = append(links, builtLink(node))

			//don't add a's children nodes to queue
			continue
		}

		//Add children nodes to queue
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			queue = append(queue, child)
		}

	}

	return links

}

//Extracts url and text from an "a href" node
func builtLink(node *html.Node) Link {
	var href string

	//Find URL
	for _, a := range node.Attr {
		if a.Key == "href" {
			href = a.Val
			break
		}
	}

	return Link{href, getLinkText(node.FirstChild)}

}

//Gets link text recursively in the following order
//1. node's text
//2. next sibling's text
//3. FirstChild's text
func getLinkText(root *html.Node) string {
	var text string

	node := root
	for node != nil {
		if node.Type == html.TextNode {
			//Remove \n and spaces to the left
			text += strings.TrimSpace(strings.ReplaceAll(node.Data, "\n", ""))
		}
		if node.FirstChild != nil {
			text += " " + getLinkText(node.FirstChild)
		}
		node = node.NextSibling
	}
	return text
}
